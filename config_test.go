package main

import (
	"testing"
)

func TestBaseConfig(t *testing.T) {
	configDef := `
    [threads]
    hash = 10
    upload = 2

    [aws]
    access = abc123Access
    secret = abc123Secret

    [backup "test"]
    region = eu-west-1
    path = /tmp/
    db = tmp.db
    exclude = *.png
    exclude = *.tiff
    include = *.jpg
    vault = test
`

	config, err := ReadConfig(configDef)
	if err != nil {
		t.Errorf("Error parsing config: %s", err)
	}

	if _, ok := config.Backup["test"]; ok == false {
		t.Errorf("Backup config for `test` not found. Skipping remaining tests for 'test' backup config.")
		return
	}

	if config.Threads.Hash != 10 {
		t.Errorf("Invalid number of hash threads `%d`, expected `%d`", config.Threads.Hash, 10)
	}

	if config.Threads.Upload != 2 {
		t.Errorf("Invalid number of upload threads `%d`, expected `%d`", config.Threads.Hash, 2)
	}

	if config.Aws.Access != "abc123Access" {
		t.Errorf("Invalid AWS access code `%s`, expected `%s`", config.Aws.Access, "abc123Access")
	}

	if config.Aws.Secret != "abc123Secret" {
		t.Errorf("Invalid AWS secret code `%s`, expected `%s`", config.Aws.Secret, "abc123Secret")
	}

	if config.Backup["test"].Region.Region.Name != "eu-west-1" {
		t.Errorf("Invalid region `%d`, expected `%d`", config.Backup["test"].Region.Region.Name, "eu-west-1")
	}

	if config.Backup["test"].Path != "/tmp/" {
		t.Errorf("Invalid path `%s`, expected `%s`", config.Backup["test"].Path, "/tmp/")
	}

	if config.Backup["test"].Db != "tmp.db" {
		t.Errorf("Invalid db `%s`, expected `%s`", config.Backup["test"].Db, "tmp.db")
	}

	if config.Backup["test"].Vault != "test" {
		t.Errorf("Invalid vault `%s`, expected `%s`", config.Backup["test"].Vault, "test")
	}

	if !compareInclusions(config.Backup["test"].Exclude, []string{"*.png", "*.tiff"}) {
		t.Errorf("Invalid exclude `%+v`, expected `%+v`", config.Backup["test"].Exclude, []string{"*.png", "*.tiff"})
	}

	if !compareInclusions(config.Backup["test"].Include, []string{"*.jpg"}) {
		t.Errorf("Invalid exclude `%+v`, expected `%+v`", config.Backup["test"].Include, []string{"*.jpg"})
	}
}

func TestNoConfigs(t *testing.T) {
	configDef := `
    [threads]
    hash = 10
    upload = 2

    [aws]
    access = 123abcAccess
    secret = 123abcSecret
`
	if _, err := ReadConfig(configDef); err == nil || err.Error() != "No configurations given" {
		t.Error("Expected error `No configurations given`")
	}
}

func TestHashThreadsMissing(t *testing.T) {
	configDef := `
    [aws]
    access = abc123Access
    secret = abc123Secret

    [backup "test"]
    region = eu-west-1
    path = /tmp/
    db = tmp.db
`

	if _, err := ReadConfig(configDef); err == nil || err.Error() != "Need at least one hash thread" {
		t.Error("ReadConfig should have complained about 0 value for hash threads")
	}
}

func TestUploadThreadsMissing(t *testing.T) {
	configDef := `
    [threads]
    hash = 10

    [aws]
    access = abc123Access
    secret = abc123Secret

    [backup "test"]
    region = eu-west-1
    path = /tmp/
    db = tmp.db
`

	if _, err := ReadConfig(configDef); err == nil || err.Error() != "Need at least one upload thread" {
		t.Error("ReadConfig should have complained about 0 value for upload threads: %s", err)
	}
}

func TestAwsConfigInBackup(t *testing.T) {
	configDef := `
    [threads]
    hash = 10
    upload = 2

    [aws]
    access = abc123Access
    secret = abc123Secret

    [backup "test"]
    region = eu-west-1
    path = /tmp/
    db = tmp.db
    vault = test
    aws-access = 123abcAccess
    aws-secret = 123abcSecret
`

	config, err := ReadConfig(configDef)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if config.Backup["test"].AwsAccess != "123abcAccess" {
		t.Errorf("Invalid AWS Access code `%s`, expected `%s`", config.Backup["test"].AwsAccess, "123abcAccess")
	}

	if config.Backup["test"].AwsSecret != "123abcSecret" {
		t.Errorf("Invalid AWS Secret `%s`, expected `%s`", config.Backup["test"].AwsSecret, "123abcSecret")
	}
}

func TestAwsConfigFromGlobal(t *testing.T) {
	configDef := `
    [threads]
    hash = 10
    upload = 2

    [aws]
    access = abc123Access
    secret = abc123Secret

    [backup "test"]
    vault = test
    region = eu-west-1
    path = /tmp/
    db = tmp.db
`
	config, err := ReadConfig(configDef)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if config.Backup["test"].AwsAccess != "abc123Access" {
		t.Errorf("Invalid AWS Access code `%s`, expected `%s`", config.Backup["test"].AwsAccess, "abc123Access")
	}

	if config.Backup["test"].AwsSecret != "abc123Secret" {
		t.Errorf("Invalid AWS Secret `%s`, expected `%s`", config.Backup["test"].AwsSecret, "abc123Secret")
	}
}

func TestNoAwsConfig(t *testing.T) {
	configDef := `
    [threads]
    hash = 10
    upload = 2

    [backup "test"]
    vault = test
    region = eu-west-1
    path = /tmp/
    db = tmp.db
`
	if _, err := ReadConfig(configDef); err == nil || err.Error() != "No AWS credentials supplied for backup `test`" {
		t.Error("Expected error about AWS credentials in `test` backup")
	}
}

func TestAwsAccessOnly(t *testing.T) {
	configDef := `
    [threads]
    hash = 10
    upload = 2

    [backup "test"]
    vault = test
    region = eu-west-1
    path = /tmp/
    db = tmp.db
    aws-access = abc123Access
`
	if _, err := ReadConfig(configDef); err == nil || err.Error() != "AWS Access code suplied, but no AWS Secret for config `test`" {
		t.Error("Expected error about AWS credentials in `test` backup.")
	}
}

func TestAwsSecretOnly(t *testing.T) {
	configDef := `
    [threads]
    hash = 10
    upload = 2

    [backup "test"]
    vault = test
    region = eu-west-1
    path = /tmp/
    db = tmp.db
    aws-secret = abc123Secret
`
	if _, err := ReadConfig(configDef); err == nil || err.Error() != "AWS Secret suplied, but no AWS Access code for config `test`" {
		t.Error("Expected error about AWS credentials in `test` backup.")
	}
}

func TestInvalidRegion(t *testing.T) {
	configDef := `
    [threads]
    hash = 10
    upload = 2

    [aws]
    access = 123abcAccess
    secret = 123abcSecret

    [backup "test"]
    region = xx-yyyy-999
    path = /tmp/
    db = tmp.db
`
	if _, err := ReadConfig(configDef); err == nil || err.Error() != "Unable to find region xx-yyyy-999" {
		t.Error("Expected error about non-existing AWS region in `test` backup")
	}
}

func TestNoRegion(t *testing.T) {
	configDef := `
    [threads]
    hash = 10
    upload = 2

    [aws]
    access = 123abcAccess
    secret = 123abcSecret

    [backup "test"]
    path = /tmp/
    db = tmp.db
`
	if _, err := ReadConfig(configDef); err == nil || err.Error() != "No region supplied for config `test`" {
		t.Error("Expected error about non-existing AWS region in `test` backup")
	}
}

func TestNoPath(t *testing.T) {
	configDef := `
    [threads]
    hash = 10
    upload = 2

    [aws]
    access = 123abcAccess
    secret = 123abcSecret

    [backup "test"]
    region = eu-west-1
    db = test.db
`
	if _, err := ReadConfig(configDef); err == nil || err.Error() != "No path supplied for config `test`" {
		t.Error("Expected error missing path param for `test` backup: %s")
	}
}

func TestNoDb(t *testing.T) {
	configDef := `
    [threads]
    hash = 10
    upload = 2

    [aws]
    access = 123abcAccess
    secret = 123abcSecret

    [backup "test"]
    region = eu-west-1
    path = /tmp/
`
	if _, err := ReadConfig(configDef); err == nil || err.Error() != "No db supplied for config `test`" {
		t.Error("Expected error missing db param for `test` backup")
	}
}

func compareInclusions(test, compare []string) bool {
	for i, val := range test {
		if compare[i] != val {
			return false
		}
	}
	return true
}
