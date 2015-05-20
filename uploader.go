package main

import (
	"errors"
	"fmt"
	"github.com/rdwilliamson/aws"
	"github.com/rdwilliamson/aws/glacier"
	"strings"
)

type Uploader struct {
	conn       *glacier.Connection
	vault      string
	indexVault string
}

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
