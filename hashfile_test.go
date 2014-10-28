package main

import (
	"testing"
)

func TestHashFile(t *testing.T) {
	files := map[string]*File{
		"32d10c7b8cf96570ca04ce37f2a19d84240d3a89": NewFile("filesets/fileset1/file1"),
		"80256f39a9d308650ac90d9be9a72a9562454574": NewFile("filesets/fileset1/file3"),
		"9888f1b80855c800640a0df40da83dfd14111456": NewFile("filesets/fileset1/sub/file1"),
		"01b307acba4f54f55aafc33bb06bbbf6ca803e9a": NewFile("filesets/fileset1/sub/file2"),
	}

	for hash, file := range files {
		HashFile(file)
		if hash != file.Hash {
			t.Errorf("Hash for `%s` was expected to be `%s`, got `%s`", file.Filename, hash, file.Hash)
		}
	}
}
