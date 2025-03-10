package source

var SupportedFileSuffixes = map[string]bool{
	"yaml": true,
	"yml":  true,
	"json": true,
	"toml": true,
}

type Source interface {
	Read() error
	String() string
	//Read() (*FileSet, error)
	//Write(*FileSet) error
	//Watch() (Watcher, error)
}
