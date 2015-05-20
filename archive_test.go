package main

import (
	"testing"
)

func TestAddAndFindByFilename(t *testing.T) {
	archive, err := NewArchive(":memory:")
	if err != nil {
		t.Errorf("Could not create archive instance: %s", err)
	}

	file := &ArchivedFile{
		filename:  "hello.txt",
		hash:      "h12345",
		amazonId:  "a12345",
		isDeleted: false,
	}

	err = archive.AddFile(file)

	if err != nil {
		t.Errorf("File should have been added, but got error: %s", err)
	}

	file, err = archive.FindFileByFilename("hello.txt")
	if err != nil {
		t.Errorf("Should be able to find file but got error: %s", err)
	} else {
		if file.Filename() != "hello.txt" {
			t.Errorf("Invalid filename `%s`, expected `hello.txt`", file.Filename())
		}
		if file.AmazonId() != "a12345" {
			t.Errorf("Invalid AmazonID `%s`, expected `a12345`", file.AmazonId())
		}
		if file.Hash() != "h12345" {
			t.Errorf("Invalid hash `%s`, expected `h12345`", file.Hash())
		}
		if file.IsDeleted() {
			t.Errorf("Expected file not be deleted, but it is")
		}
	}
}

func TestAddAndFindByHash(t *testing.T) {
	archive, err := NewArchive(":memory:")
	if err != nil {
		t.Errorf("Could not create archive instance: %s", err)
	}

	file := &ArchivedFile{
		filename:  "hello.txt",
		hash:      "h12345",
		amazonId:  "a12345",
		isDeleted: false,
	}

	err = archive.AddFile(file)

	if err != nil {
		t.Errorf("File should have been added, but got error: %s", err)
	}

	amazonId, err := archive.FindAmazonIdByHash("h12345")
	if err != nil {
		t.Errorf("Should be able to find file but got error: %s", err)
	} else {
		if *amazonId != "a12345" {
			t.Errorf("Invalid AmazonID `%s`, expected `a12345`", amazonId)
		}
	}
}

func TestAddAndFindInEmptyArchive(t *testing.T) {
	archive, err := NewArchive(":memory:")
	if err != nil {
		t.Errorf("Could not create archive instance: %s", err)
	}

	_, err = archive.FindFileByFilename("foo.txt")
	if err == nil {
		t.Errorf("Expected error when searching for file in empty archive but got no error.")
	}

	_, err = archive.FindAmazonIdByHash("a12345")
	if err == nil {
		t.Errorf("Expected error when searching for file in empty archive but got no error.")
	}
}

func TestAddAndDeleteFile(t *testing.T) {
	archive, err := NewArchive(":memory:")
	if err != nil {
		t.Errorf("Could not create archive instance: %s", err)
	}

	file := &ArchivedFile{
		filename:  "hello.txt",
		hash:      "h12345",
		amazonId:  "a12345",
		isDeleted: false,
	}

	err = archive.AddFile(file)

	if err != nil {
		t.Errorf("File should have been added, but got error: %s", err)
	}

	err = archive.DeleteFile(file.Hash(), file.Filename())
	if err != nil {
		t.Errorf("Error deleting file: %s", err)
	}

	_, err = archive.FindFileByFilename("hello.txt")
	if err == nil {
		t.Errorf("Should not be able to find file since it was deleted")
	}

	_, err = archive.FindAmazonIdByHash("h12345")
	if err != nil {
		t.Errorf("AmazonID for hash `h12345` should still be in DB even though file is deleted")
	}
}
