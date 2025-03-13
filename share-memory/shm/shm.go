package shm

import (
	"fmt"
	"io"
	"unsafe"

	"github.com/pkg/errors"
	"golang.org/x/sys/unix"

	"github.com/peanut-io/peanut/logger"
)

const (
	R  = 0400
	W  = 0200
	RW = R | W
)

const (
	Key  = 0x25312
	Size = 1 << 31
)

type Shm struct {
	Id, addr, key, size, offset int
}

func NewShm(id, size, key, perm int) (*Shm, error) {
	m := &Shm{
		size: ((size | -size) & size) | (^(size | -size) & Size),
		key:  ((key | -key) & key) | (^(key | -key) & Key),
		Id:   id,
	}
	if err := m.initialize(perm); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *Shm) initialize(perm int) error {
	if s.Id != 0 {
		logger.Warnw("share memory intend to bound process")
		return s.attach()
	}
	shmId, _, err := unix.Syscall(unix.SYS_SHMGET, uintptr(s.key), uintptr(s.size), (unix.IPC_CREAT|unix.IPC_EXCL)|uintptr(perm))
	if err == unix.EEXIST && int(shmId) < 0 {
		shmId, _, err = unix.Syscall(unix.SYS_SHMGET, uintptr(s.key), 0, 0666)
	}
	if err != 0 {
		logger.Warnw("share memory failed to create", "error", err.Error())
		return errors.Wrap(err, fmt.Sprintf("share memory failed to create, error %s ", err.Error()))
	}
	s.Id = int(shmId)
	return s.attach()
}

func (s *Shm) Write(buffer []byte) (int, error) {
	var offset, length int

	if s.offset >= s.size {
		return 0, io.EOF
	}

	length = len(buffer)
	if length+s.offset > s.size {
		length = s.size - s.offset
	}

	for ; offset < length; offset++ {
		addr := (*byte)(unsafe.Pointer(uintptr(s.addr + s.offset + offset)))
		*addr = buffer[offset]
	}

	s.offset += offset
	return offset, io.EOF
}

func (s *Shm) Read(buffer []byte) (int, error) {
	var offset, length int

	if s.offset >= s.size {
		return 0, io.EOF
	}

	length = len(buffer)
	if length+s.offset > s.size {
		length = s.size - s.offset
	}

	for ; offset < length; offset++ {
		b := (*byte)(unsafe.Pointer(uintptr(s.addr + s.offset + offset)))
		buffer[offset] = *b
	}

	s.offset += offset
	return offset, io.EOF
}

func (s *Shm) Seek(offset int, whence int) (int, error) {
	var newOffset int
	switch whence {
	case 1:
		newOffset = s.offset + offset
	case 2:
		newOffset = s.size - offset
	default:
		newOffset = offset
	}
	if newOffset < 0 {
		return 0, errors.New("cannot seek to position before start of segment")
	}
	s.offset = newOffset
	return newOffset, nil
}

// Close delete memory
func (s *Shm) Close() error {
	if s.addr != 0 {
		_ = s.detach()
		_, _, err := unix.Syscall(unix.SYS_SHMCTL, uintptr(s.Id), unix.IPC_RMID, 0)
		if err != 0 && err != unix.EINVAL {
			logger.Warnw("failed to delete the share memory", "error", err.Error())
			return errors.Wrap(err, fmt.Sprintf("shmctl failed to delete the share memory, error %s ", err))
		}
	}
	return nil
}

// attach map memory
func (s *Shm) attach() error {
	addr, _, err := unix.Syscall(unix.SYS_SHMAT, uintptr(s.Id), 0, 0)
	if int(addr) == -1 {
		logger.Warnw("share memory failed to map process", "error", err.Error())
		return errors.Wrap(err, fmt.Sprintf("shmat failed to map process, error %s ", err.Error()))
	}
	s.addr = int(addr)
	return nil
}

// detach unmap memory
func (s *Shm) detach() error {
	_, _, err := unix.Syscall(unix.SYS_SHMDT, uintptr(s.addr), 0, 0)
	if err != 0 {
		logger.Warnw("share memory failed to unmap process", "error", err.Error())
		return errors.Wrap(err, fmt.Sprintf("shmdt failed to unmap process, error %s ", err.Error()))
	}
	return nil
}

func (s *Shm) Clear() {
	s.memorySet(s.offset)
	s.offset = 0
}

func (s *Shm) memorySet(offset int) {
	var i int
	for ; i < offset; i++ {
		*(*byte)(unsafe.Pointer(uintptr(s.addr + i))) = 0
	}
}

func (s *Shm) IsEmpty() bool {
	return s.addr == 0 || s.Id == 0
}

func (s *Shm) String() string {
	return fmt.Sprintf("id: %d ,size: %d, key: %d", s.Id, s.size, s.key)
}

func (s *Shm) State() (*unix.SysvShmDesc, error) {
	shmStat := &unix.SysvShmDesc{}
	pointer := unsafe.Pointer(shmStat)
	_, _, err := unix.Syscall(unix.SYS_SHMCTL, uintptr(s.Id), unix.IPC_STAT, uintptr(pointer))
	if err != 0 {
		return nil, err
	}
	return shmStat, nil
}
