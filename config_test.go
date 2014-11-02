package main

import (
	"testing"
)

func TestSimpleConfig(t *testing.T) {
	configDef := `
    [threads]
    hash = 10
    upload = 2

    [aws]
    access = abc123
    secret = abc123

    [backup "test"]
    region = eu-west-1
    path = /tmp/
    db = tmp.db
    exclude = *.png,*.tiff
    include = *.jpg
`

	config, err := ReadConfig(configDef)
	if err != nil {
		t.Errorf("Error parsing config: %s", err)
	}

	if _, ok := config.Backup["test"]; ok == false {
		t.Errorf("Backup config for `test` not found. Skipping remaining tests for 'test' backup config.")
		return
	}

	if config.Backup["test"].Path != "/tmp/" {
		t.Errorf("Invalid path `%s`, expected `%s`", config.Backup["test"].Path, "/tmp/")
	}

	if config.Backup["test"].Db != "tmp.db" {
		t.Errorf("Invalid db `%s`, expected `%s`", config.Backup["test"].Db, "tmp.db")
	}

	if config.Backup["test"].Exclude != "*.png,*.tiff" {
		t.Errorf("Invalid exclude `%s`, expected `%s`", config.Backup["test"].Exclude, "*.png,*.tiff")
	}

	if config.Backup["test"].Include != "*.jpg" {
		t.Errorf("Invalid exclude `%s`, expected `%s`", config.Backup["test"].Include, "*.jpg")
	}
}

func TestThreadsMissing(t *testing.T) {
	configDef := `
    [aws]
    access = abc123
    secret = abc123

    [backup "test"]
    region = eu-west-1
    path = /tmp/
    db = tmp.db
    exclude = *.png,*.tiff
    include = *.jpg
`

	if _, err := ReadConfig(configDef); err == nil {
		t.Error("ReadConfig should have complained about missing values for `Threads` section")
	}
}
