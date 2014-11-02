package main

import (
	"code.google.com/p/gcfg"
	"fmt"
	"github.com/rdwilliamson/aws"
)

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
		Enabled   *bool
		Region    MyAwsRegion
		Path      string
		Db        string
		Exclude   string
		Include   string
		AwsAccess *string
		AwsSecret *string
	}
}

type MyAwsRegion struct {
	*aws.Region
}

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

func ReadConfig(configDef string) (*Config, error) {
	cfg := Config{}
	err := gcfg.ReadStringInto(&cfg, configDef)

	if cfg.Threads.Hash < 1 {
		return nil, fmt.Errorf("Need at least one hash thread")
	}
	if cfg.Threads.Upload < 1 {
		return nil, fmt.Errorf("Need at least one upload thread")
	}

	return &cfg, err
}
