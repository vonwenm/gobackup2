package main

import (
	"testing"
)

func TestHashFile(t *testing.T) {
	files := map[string]*File{
		"32d10c7b8cf96570ca04ce37f2a19d84240d3a89": NewFile("filesets/fileset1/file1.txt"),
		"80256f39a9d308650ac90d9be9a72a9562454574": NewFile("filesets/fileset1/file3.txt"),
		"9888f1b80855c800640a0df40da83dfd14111456": NewFile("filesets/fileset1/sub/file1.bin"),
		"01b307acba4f54f55aafc33bb06bbbf6ca803e9a": NewFile("filesets/fileset1/sub/file2.txt"),
	}

	for expectedHash, file := range files {
		actualHash, err := file.Hash()
		if err != nil {
			t.Errorf("Unexpected error when calculating hash for `%s`: %s", file.Filename(), err)
		}
		if actualHash != expectedHash {
			t.Errorf("Hash for `%s` was expected to be `%s`, got `%s`", file.Filename(), expectedHash, actualHash)
		}
	}
}
