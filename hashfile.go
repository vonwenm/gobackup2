package main

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
)

func HashFile(f *File) error {
	rawReader, err := os.Open(f.Filename)
	defer rawReader.Close()
	if err != nil {
		return err
	}

	reader := bufio.NewReaderSize(rawReader, 1024*1024)
	hasher := sha1.New()
	_, err = io.Copy(hasher, reader)
	if err != nil {
		return err
	}
	f.Hash = string(fmt.Sprintf("%x", hasher.Sum(nil)))
	return nil
}
