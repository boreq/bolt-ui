package frontend

import (
	"embed"
	"net/http"
	"path"
)

//go:embed static/*
var content embed.FS

type FrontendFileSystem struct {
	fs http.FileSystem
}

func NewFrontendFileSystem() (*FrontendFileSystem, error) {
	return &FrontendFileSystem{
		fs: http.FS(content),
	}, nil
}

func (f *FrontendFileSystem) Open(name string) (http.File, error) {
	file, err := f.fs.Open(f.addPrefix(name))
	if err != nil {
		file, err := f.fs.Open(f.addPrefix("/index.html"))
		if err != nil {
			return nil, err
		}
		return file, nil
	}
	return file, nil
}

func (f *FrontendFileSystem) addPrefix(name string) string {
	return path.Join("static", name)
}
