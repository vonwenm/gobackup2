package main

import (
	"io"
	"os"
)

type File struct {
	Filename string
	Hash     string
}

func (f *File) GetReader() (io.Reader, error) {
	return os.Open(f.Filename)
}
