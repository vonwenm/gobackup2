package main

import "testing"

func TestListAllFiles(t *testing.T) {
	expected := []*File{
		NewFile("filesets/fileset1/file1.txt"),
		NewFile("filesets/fileset1/file3.txt"),
		NewFile("filesets/fileset1/sub/file1.bin"),
		NewFile("filesets/fileset1/sub/file2.txt"),
	}

	c := ListFiles("./filesets/fileset1", []string{}, []string{})
	checkFiles(t, c, expected)
}

func TestListIncludeFiles(t *testing.T) {
	expected := []*File{
		NewFile("filesets/fileset1/file1.txt"),
		NewFile("filesets/fileset1/file3.txt"),
		NewFile("filesets/fileset1/sub/file2.txt"),
	}

	c := ListFiles("./filesets/fileset1", []string{"*.txt"}, []string{})
	checkFiles(t, c, expected)
}

func TestListExcludeFiles(t *testing.T) {
	expected := []*File{
		NewFile("filesets/fileset1/file3.txt"),
		NewFile("filesets/fileset1/sub/file2.txt"),
	}

	c := ListFiles("./filesets/fileset1", []string{}, []string{"*file1*"})
	checkFiles(t, c, expected)
}

func TestListIncludeAndExcludeFiles(t *testing.T) {
	expected := []*File{
		NewFile("filesets/fileset1/file1.txt"),
	}

	c := ListFiles("./filesets/fileset1", []string{"*file1*"}, []string{"*.bin"})
	checkFiles(t, c, expected)
}

func checkFiles(t *testing.T, c <-chan *File, expected []*File) {
	actual := []*File{}
	for file := range c {
		actual = append(actual, file)
		if !findInFiles(file, expected) {
			t.Errorf("Unexepected file `%s` listed!", file.Filename)
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
