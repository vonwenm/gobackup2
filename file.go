package main

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
)

/**
 * File represents a file on the filesystem
 */
type File struct {
	filename string
	hash     string
}

/**
 * NewFile creates a new File instance
 * @param string filename The filename
 */
func NewFile(filename string) *File {
	return &File{
		filename: filename,
	}
}

/**
 * Filename returns the filename of this file
 * @return string
 */
func (f *File) Filename() string {
	return f.filename
}

/**
 * Hash calculates the SHA1-hash of the file
 * and caches it. Any consequetive call of Hash
 * will return the cached value.
 * @return string The SHA1-hash of the file
 */
func (f *File) Hash() (string, error) {
	if f.hash != "" {
		return f.hash, nil
	}
	rawReader, err := os.Open(f.filename)
	defer rawReader.Close()
	if err != nil {
		return "", err
	}

	reader := bufio.NewReaderSize(rawReader, 1024*1024)
	hasher := sha1.New()
	_, err = io.Copy(hasher, reader)
	if err != nil {
		return "", err
	}
	f.hash = string(fmt.Sprintf("%x", hasher.Sum(nil)))
	return f.hash, nil
}
