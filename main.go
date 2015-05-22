package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"sync"
)

var wg sync.WaitGroup

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	log.Printf("Using %d cores\n", runtime.NumCPU())
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

	if config.Threads.Hash > runtime.NumCPU() {
		log.Printf("You want to use %d threads for hashing, but you only have %d cores available.", config.Threads.Hash, runtime.NumCPU())
		log.Printf("Even though this will work just fine, using %d hash threads is likely to give better throughput.", runtime.NumCPU())
		log.Printf("Note that for typical hard disks hashing is I/O bound, not CPU bound.")
	}

	for _, backup := range config.Backup {
		uploader, err := NewUploader(backup.AwsSecret, backup.AwsAccess, backup.Region.Region, backup.Vault)
		if err != nil {
			log.Printf("Error creating uploader: %s", err)
			continue
		}

		archive, err := NewArchive(backup.Db)
		if err != nil {
			log.Printf("Error creating archive: %s", err)
			continue
		}

		files, err := archive.ListFiles()
		for _, file := range files {
			info, err := os.Stat(file.Filename())
			if err != nil || info.IsDir() {
				archive.DeleteFile(file.Hash(), file.Filename())
			}
		}

		_, err = NewFileChecker(archive)
		if err != nil {
			log.Printf("Unable to start file checker: %s", err)
		}

		filesChan := make(chan *File, 100)
		uploadsChan := make(chan *File, 100)
		for i := 0; i < config.Threads.Hash; i++ {
			wg.Add(1)
			go Hash(filesChan, uploadsChan)
		}
		for i := 0; i < config.Threads.Upload; i++ {
			wg.Add(1)
			go Upload(uploader, uploadsChan)
		}
		ListFiles(backup.Path, backup.Include, backup.Exclude, filesChan)
	}
	wg.Wait()
}

func Hash(files chan *File, uploads chan *File) {
	defer wg.Done()
	for {
		file, ok := <-files
		if !ok {
			return
		}
		hash, err := file.Hash()
		if err != nil {
			log.Printf("Could not calculate hash for %s: %s", file.Filename(), err)
			continue
		}
		log.Printf("File: %s, Hash: %s\n", file.Filename(), hash)
	}
}

func Upload(uploader *Uploader, uploads chan *File) {
	defer wg.Done()
	for {
		file, ok := <-uploads
		if !ok {
			return
		}
		uploader.UploadFile(file.Filename())
	}
}
