package mmap

import (
	"io"
	"os"
	"unsafe"

	"github.com/pkg/errors"
	"golang.org/x/sys/unix"

	"github.com/peanut-io/peanut/logger"
)

const (
	FileName = "/tmp/mmap"
	Size     = 1 << 31
)

type MMap struct {
	fd       uintptr
	buf      unsafe.Pointer
	size     int
	offset   int
	fileName string
}

func NewMMap(fileName string, size int) (*MMap, error) {
	if len(fileName) == 0 {
		fileName = FileName
	}
	m := &MMap{
		size:     ((size | -size) & size) | (^(size | -size) & Size),
		fileName: fileName,
	}
	if err := m.initialize(); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *MMap) initialize() error {
	var buf []byte

	file, err := os.OpenFile(m.fileName, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		logger.Warnw("share memory failed to open the file", "fileName", m.fileName, "error", err.Error())
		return err
	}
	stat, err := os.Stat(m.fileName)
	if err != nil {
		return err
	}
	if stat.Size() < int64(m.size) {
		err = os.Truncate(m.fileName, int64(m.size))
		if err != nil {
			logger.Warnw("share memory failed to truncate", "error", err.Error())
			return err
		}
	}

	buf, err = unix.Mmap(int(file.Fd()), 0, int(m.size), unix.PROT_READ|unix.PROT_WRITE, unix.MAP_SHARED)
	if err != nil {
		logger.Warnw("share memory failed to map", "error", err.Error())
		return err
	}
	//m.fd = uintptr(unsafe.Pointer(file))
	m.fd = file.Fd()
	m.buf = unsafe.Pointer(&buf)
	return nil
}

func (m *MMap) Write(buffer []byte) (int, error) {
	var length int

	if m.offset >= m.size {
		return 0, errors.New("mmap buffer overflow")
	}

	length = len(buffer)
	if length+m.offset > m.size {
		length = m.size - m.offset
	}

	buf := *(*[]byte)(m.buf)
	copy(buf[m.offset:], buffer[:length])

	m.offset += length
	return length, nil
}

func (m *MMap) Read(buffer []byte) (int, error) {
	var length int

	if m.offset >= m.size {
		return 0, io.EOF
	}

	length = len(buffer)
	if length+m.offset > m.size {
		length = m.size - m.offset
	}
	buf := *(*[]byte)(m.buf)
	copy(buffer, buf[m.offset:m.offset+length])
	m.offset += length
	return length, nil
}

func (m *MMap) Close() (err error) {
	if m.fd != 0 {
		file := os.NewFile(m.fd, m.fileName)
		data := *(*[]byte)(m.buf)
		err = unix.Munmap(data)
		err = file.Close()
	}
	return
}

func (m *MMap) Seek(offset int, whence int) (int, error) {
	var newOffset int
	switch whence {
	case 1:
		newOffset = m.offset + offset
	case 2:
		newOffset = m.size - offset
	default:
		newOffset = offset
	}
	if newOffset < 0 {
		return 0, errors.New("Cannot seek to position before start of segment")
	}
	m.offset = newOffset
	return newOffset, nil
}

func (m *MMap) Clear() {
	m.memorySet(m.offset)
	m.offset = 0
}

func (m *MMap) memorySet(offset int) {
	copy(*(*[]byte)(m.buf), make([]byte, offset))
}
