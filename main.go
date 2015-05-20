package main

import (
	"flag"
	"io/ioutil"
	"log"
)

func main() {
	configFile := flag.String("config", "gobackup.ini", "Path to config file")
	flag.Parse()

	configDef, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	config, err := ReadConfig(string(configDef))
	if err != nil {
		log.Fatalf("Error parsing config: %s", err)
	}

	for _, backup := range config.Backup {
		_, err := NewUploader(backup.AwsSecret, backup.AwsAccess, backup.Region.Region, backup.Vault)
		if err != nil {
			log.Printf("Error creating uploader: %s", err)
			continue
		}

		_, err = NewArchive(backup.Db)
		if err != nil {
			log.Printf("Error creating archive: %s", err)
		}
	}
}
