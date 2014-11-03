package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

/**
 * List all files in a given path recursively
 * Omits all directories and empty files
 */
func ListFiles(path string, include, exclude []string) <-chan *File {
	ch := make(chan *File)

	inRegex := regexify(include, ".*")
	exRegex := regexify(exclude, "")

	go func() {
		filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() || info.Size() == 0 {
				return nil
			}

			if inRegex != nil && !inRegex.Match([]byte(path)) {
				return nil
			}

			if exRegex != nil && exRegex.Match([]byte(path)) {
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

func regexify(patterns []string, ifEmpty string) *regexp.Regexp {
	var str string
	if len(patterns) == 0 {
		return nil
	} else {
		for i := range patterns {
			patterns[i] = strings.Replace(patterns[i], ".", "\\.", -1)
			patterns[i] = strings.Replace(patterns[i], "*", ".*", -1)
			patterns[i] = "^" + patterns[i] + "$"
		}
		str = strings.Join(patterns, "|")
	}
	return regexp.MustCompile(str)
}
