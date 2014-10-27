package main

import (
	"bytes"
	"io"
)

type Testfile struct {
	Filename string
	Data     string
	Hash     string
}

func (t *Testfile) GetReader() (io.Reader, error) {
	return bytes.NewBufferString(t.Data), nil
}
