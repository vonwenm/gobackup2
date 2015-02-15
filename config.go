package main

import (
	"code.google.com/p/gcfg"
	"errors"
	"fmt"
	"github.com/rdwilliamson/aws"
)

/**
 * Main Config struct for program
 */
type Config struct {
	Threads struct {
		Hash   int
		Upload int
	}
	Aws struct {
		Secret string
		Access string
	}
	Backup map[string]*struct {
		Region    MyAwsRegion
		Path      string
		Db        string
		Exclude   []string
		Include   []string
		AwsAccess string `gcfg:"aws-access"`
		AwsSecret string `gcfg:"aws-secret"`
	}
}

/**
 * Simple wrapper for aws.Region
 * Allowing us to add a custom unmarshal method
 */
type MyAwsRegion struct {
	*aws.Region
}

/**
 * Custom unmarshaller for MyAwsRegion
 * Look up region in aws.Regions map
 * @return error Returns error if region does not exist in aws.Regions
 */
func (m *MyAwsRegion) UnmarshalText(text []byte) error {
	region := string(text)
	for _, r := range aws.Regions {
		if r.Name == region {
			*m = MyAwsRegion{r}
			return nil
		}
	}
	return fmt.Errorf("Unable to find region %s", region)
}

/**
 * Read config from a string
 * @param string The configuration definition
 * @return Config A Config struct, or nil if something went wrong
 * @return error Returns error if something went wrong. Config will be nil in this case.
 */
func ReadConfig(configDef string) (*Config, error) {
	cfg := Config{}
	if err := gcfg.ReadStringInto(&cfg, configDef); err != nil {
		return nil, err
	}

	if cfg.Threads.Hash < 1 {
		return nil, fmt.Errorf("Need at least one hash thread")
	}
	if cfg.Threads.Upload < 1 {
		return nil, fmt.Errorf("Need at least one upload thread")
	}

	if len(cfg.Backup) == 0 {
		return nil, errors.New("No configurations given")
	}

	for key, backup := range cfg.Backup {
		if backup.Region.Region == nil {
			return nil, fmt.Errorf("No region supplied for config `%s`", key)
		}

		if backup.Path == "" {
			return nil, fmt.Errorf("No path supplied for config `%s`", key)
		}

		if backup.Db == "" {
			return nil, fmt.Errorf("No db supplied for config `%s`", key)
		}

		if backup.AwsAccess != "" && backup.AwsSecret == "" {
			return nil, fmt.Errorf("AWS Access code suplied, but no AWS Secret for config `%s`", key)
		}

		if backup.AwsAccess == "" && backup.AwsSecret != "" {
			return nil, fmt.Errorf("AWS Secret suplied, but no AWS Access code for config `%s`", key)
		}

		if backup.AwsAccess == "" && backup.AwsSecret == "" {
			if cfg.Aws.Access != "" && cfg.Aws.Secret != "" {
				backup.AwsAccess = cfg.Aws.Access
				backup.AwsSecret = cfg.Aws.Secret
			} else {
				return nil, fmt.Errorf("No AWS credentials supplied for backup `%s`", key)
			}
		}
	}

	return &cfg, nil
}
