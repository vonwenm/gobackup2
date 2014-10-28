package main

import "testing"

func TestListFiles(t *testing.T) {
	expected := []*File{
		NewFile("filesets/fileset1/file1"),
		NewFile("filesets/fileset1/file3"),
		NewFile("filesets/fileset1/sub/file1"),
		NewFile("filesets/fileset1/sub/file2"),
	}
	actual := []*File{}

	c := ListFiles("./filesets/fileset1")
	for file := range c {
		actual = append(actual, file)
		if !findInFiles(file, expected) {
			t.Errorf("Could not find `%s` in expected files!", file.Filename)
		}
	}

	if len(expected) != len(actual) {
		t.Errorf("Expected %d files, but found %d", len(expected), len(actual))
	}
}

func findInFiles(file *File, expected []*File) bool {
	for _, file2 := range expected {
		if file.Filename == file2.Filename {
			return true
		}
	}
	return false
}
