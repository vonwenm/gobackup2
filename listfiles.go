package main

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

/**
 * ListFiles lists all files in a given path recursively.
 * It omits directories and empty files. When passed a slice of includes,
 * filenames must match one of the includes to be returned.
 * When passed a slice of excludes filenames must
 * not match any of the excludes, otherwise they won't be return.
 * If a filename matches both include and exclude, it will be excluded.
 * This function closes the channel when it's done looping all files.
 * @param path string The path to scan
 * @param include []string Slice of include patterns
 * @param exclude []string Slice of exclude patterns
 * @param out <-chan *File
 */
func ListFiles(path string, include, exclude []string, out chan<- *File) {
	inRegex, exRegex := regexify(include), regexify(exclude)

	go func() {
		filepath.Walk(path, func(path string, info os.FileInfo, err error) (outErr error) {
			if err != nil {
				if info.IsDir() {
					log.Printf("Error reading directory %s: %s. Skipping.", path, err)
					return filepath.SkipDir
				} else {
					log.Printf("Error reading file %s: %s. Skipping file.", path, err)
					return
				}
			}

			if info.IsDir() || info.Size() == 0 {
				return
			}

			if inRegex != nil && !inRegex.Match([]byte(path)) {
				return
			}

			if exRegex != nil && exRegex.Match([]byte(path)) {
				return
			}

			file := NewFile(path)
			out <- file
			return
		})
		close(out)
	}()
}

/**
 * regexify changes a list of human readable patterns to
 * a regular expression. The regular expression is case insensitive
 * to make up for different capitalisation in filenames.
 * i.e. ["*.jpg", "*.png"] becomes ^.*\.jpg$|^.*\.png$
 * @param patterns []string Input patterns
 * @return
 */
func regexify(patterns []string) *regexp.Regexp {
	var str string
	if len(patterns) == 0 {
		return nil
	}

	// first copy patterns so we don't change
	// them outside our function
	p := make([]string, len(patterns))
	copy(p, patterns)

	// now loop over the patterns and turn them into a
	// (sub) regular expressions
	for i := range p {
		p[i] = strings.Replace(p[i], ".", "\\.", -1)
		p[i] = strings.Replace(p[i], "*", ".*", -1)
	}

	// now glue the regular expression parts together
	str = "(?i)^" + strings.Join(p, "|") + "$"
	return regexp.MustCompile(str)
}
