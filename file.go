package main

import (
	"io"
	"os"
)

type File struct {
	Filename string
	Hash     string
}

func NewFile(filename string) File {
	return File{
		Filename: filename,
	}
}

func (f *File) GetReader() (io.Reader, error) {
	return os.Open(f.Filename)
}
