package main

import (
	"os"
	"path/filepath"
)

/**
 * List all files in a given path recursively
 * Omits all directories and empty files
 */
func ListFiles(path string) <-chan File {
	ch := make(chan File)
	go func() {
		filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() || info.Size() == 0 {
				return nil
			}
			file := NewFile(path)
			ch <- file
			return nil
		})
		close(ch)
	}()
	return ch
}
