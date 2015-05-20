package main

import (
	"errors"
	"fmt"
	"github.com/rdwilliamson/aws"
	"github.com/rdwilliamson/aws/glacier"
	"os"
	"strings"
)

/**
 * Uploader is responsible for uploading files to AWS Glacier
 */
type Uploader struct {
	conn       *glacier.Connection
	vault      string
	indexVault string
}

/**
 * NewUploader creates a new uploader instance
 * and makes sure all needed vaults exist. If they don't exist
 * they will be created
 */
func NewUploader(awsSecret, awsAccess string, awsRegion *aws.Region, vault string) (*Uploader, error) {
	if strings.HasSuffix(vault, "_index") {
		return nil, errors.New("Vault names can not end in `_index`")
	}

	conn := glacier.NewConnection(awsSecret, awsAccess, awsRegion)
	indexVault := vault + "_index"

	vaults, _, err := conn.ListVaults("", 0)
	if err != nil {
		return nil, err
	}

	mainVaultFound, indexVaultFound := false, false
	for _, i := range vaults {
		if vault == i.VaultName {
			mainVaultFound = true
		}
		if indexVault == i.VaultName {
			indexVaultFound = true
		}
	}

	if !mainVaultFound {
		if err := conn.CreateVault(vault); err != nil {
			return nil, fmt.Errorf("Vault `%s` does not exist and could not be created", vault)
		}
	}

	if !indexVaultFound {
		if err := conn.CreateVault(indexVault); err != nil {
			return nil, fmt.Errorf("Vault `%s` does not exist and could not be created", indexVault)
		}
	}

	return &Uploader{
		conn:       conn,
		vault:      vault,
		indexVault: indexVault,
	}, nil
}

/**
 * UploadFile tries to upload a file to AWS glacier.
 * Will bail after 3 failed attempts.
 */
func (u *Uploader) UploadFile(path string) (amazonId string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	for retries := 1; retries <= 3; retries++ {
		f.Seek(0, 0)
		if amazonId, err = u.conn.UploadArchive(u.vault, f, path); err != nil {
			if retries == 3 {
				err = fmt.Errorf("Upload failed after 3 retries: %s", err)
				return
			}
		} else {
			return
		}
	}
	return
}
