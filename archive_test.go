package main

import (
	"reflect"
	"testing"
)

func TestAndListMulti(t *testing.T) {
	archive, err := NewArchive(":memory:")
	if err != nil {
		t.Errorf("Could not create archive instance: %s", err)
	}

	files := []*ArchivedFile{
		&ArchivedFile{
			filename:  "hello.txt",
			hash:      "h12345",
			amazonId:  "a12345",
			isDeleted: false,
		},
		&ArchivedFile{
			filename:  "bye.txt",
			hash:      "h54321",
			amazonId:  "a54321",
			isDeleted: false,
		},
	}

	for _, file := range files {
		err = archive.AddFile(file)

		if err != nil {
			t.Errorf("File should have been added, but got error: %s", err)
		}
	}

	listedFiles, err := archive.ListFiles()
	if err != nil {
		t.Errorf("Unexpected error while listing files: %s", err)
	}

	if len(listedFiles) != len(files) {
		t.Errorf("Expected %d file, got %d", len(files), len(listedFiles))
	}

	for _, file := range files {
		if !findFileInList(file, listedFiles) {
			t.Errorf("Could not find file `%s` in listed files", file.Filename())
		}
	}
}

func TestAndListSameHash(t *testing.T) {
	archive, err := NewArchive(":memory:")
	if err != nil {
		t.Errorf("Could not create archive instance: %s", err)
	}

	files := []*ArchivedFile{
		&ArchivedFile{
			filename:  "hello.txt",
			hash:      "h12345",
			amazonId:  "a12345",
			isDeleted: false,
		},
		&ArchivedFile{
			filename:  "bye.txt",
			hash:      "h12345",
			amazonId:  "a12345",
			isDeleted: false,
		},
	}

	for _, file := range files {
		err = archive.AddFile(file)

		if err != nil {
			t.Errorf("File should have been added, but got error: %s", err)
		}
	}

	conn := archive.Connection()
	stmt, err := conn.Prepare("SELECT count(*) FROM upload")
	if err != nil {
		t.Errorf("Unexpected error when executing query: %s", err)
	}

	row := stmt.QueryRow()
	var count int
	row.Scan(&count)

	if count != 1 {
		t.Errorf("Expected 1 row in `upload` table, found %d", count)
	}

	stmt, err = conn.Prepare("SELECT count(*) FROM file")
	if err != nil {
		t.Errorf("Unexpected error when executing query: %s", err)
	}

	row = stmt.QueryRow()
	row.Scan(&count)

	if count != 2 {
		t.Errorf("Expected 2 rows in `file` table, found %d", count)
	}
}

func findFileInList(file *ArchivedFile, files []*ArchivedFile) bool {
	for _, compare := range files {
		if reflect.DeepEqual(file, compare) {
			return true
		}
	}
	return false
}

func TestAddAndFindByFilename(t *testing.T) {
	archive, err := NewArchive(":memory:")
	if err != nil {
		t.Errorf("Could not create archive instance: %s", err)
	}

	originalFile := &ArchivedFile{
		filename:  "hello.txt",
		hash:      "h12345",
		amazonId:  "a12345",
		isDeleted: false,
	}

	err = archive.AddFile(originalFile)

	if err != nil {
		t.Errorf("File should have been added, but got error: %s", err)
	}

	file, err := archive.FindFileByFilename("hello.txt")
	if err != nil {
		t.Errorf("Should be able to find file but got error: %s", err)
	} else {
		if !reflect.DeepEqual(file, originalFile) {
			t.Errorf("File retrieved fia FindFileByFilename is not the same as the one that was added.")
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

	listedFiles, err := archive.ListFiles()
	if err != nil {
		t.Errorf("Unexpected error while listing files: %s", err)
	}

	if len(listedFiles) != 0 {
		t.Errorf("Expected no files in list, got %d", len(listedFiles))
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

	listedFiles, err := archive.ListFiles()
	if err != nil {
		t.Errorf("Unexpected error while listing files: %s", err)
	}

	if len(listedFiles) != 0 {
		t.Errorf("Expected no files in list, got %d", len(listedFiles))
	}
}
