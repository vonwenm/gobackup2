package main

type File struct {
	Filename string
	Hash     string
}

func NewFile(filename string) *File {
	return &File{
		Filename: filename,
	}
}
