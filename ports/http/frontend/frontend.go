package frontend

import (
	"embed"
	"net/http"

	"github.com/boreq/velo/logging"
)

//go:embed css/* js/* index.html favicon.ico
var content embed.FS

type FrontendFileSystem struct {
	fs  http.FileSystem
	log logging.Logger
}

func NewFrontendFileSystem() (*FrontendFileSystem, error) {
	return &FrontendFileSystem{
		fs:  http.FS(content),
		log: logging.New("frontend"),
	}, nil
}

func (f *FrontendFileSystem) Open(name string) (http.File, error) {
	f.log.Debug("serving frontend file", "name", name)

	file, err := f.fs.Open(name)
	if err != nil {
		file, err := f.fs.Open("/index.html")
		if err != nil {
			return nil, err
		}
		return file, nil
	}
	return file, nil
}
