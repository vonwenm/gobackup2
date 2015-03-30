package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

/**
 * ListFiles lists all files in a given path recursively.
 * It omits and empty files. When passed a slice of includes,
 * filenames must match one of the includes to be returned.
 * When passed a slice of excludes filenames must
 * not match any of the excludes, otherwise they won't be return
 * If a filename matches both include and exclude, it will be excluded.
 * @param path string The path to scan
 * @param include []string Slice of include patterns
 * @param exclude []string Slice of exclude patterns
 * @param out <-chan *File
 */
func ListFiles(path string, include, exclude []string, out chan<- *File) {
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
			out <- file
			return nil
		})
		close(out)
	}()
}

/**
 * Regexify changes a list of human readable patterns to
 * a regular expression. The regular expression is case insensitive
 * to make up for different capitalisation in filenames.
 * i.e. ["*.jpg", "*.png"] becomes ^.*\.jpg$|^.*\.png$
 * @param patterns []string Input patterns
 * @param ifEmpty string The default input pattern when the patterns slice is empty
 * @return
 */
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
		str = "(?i)" + strings.Join(patterns, "|")
	}
	return regexp.MustCompile(str)
}
