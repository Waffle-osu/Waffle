package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

func RenameOszs(oszDir string) {
	//Setup Logger
	filename := fmt.Sprintf("logs/%d-log-osz_rename.txt", time.Now().Unix())

	file, fileErr := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if fileErr != nil {
		panic(fileErr)
	}

	multiWriter := io.MultiWriter(file, os.Stdout)

	logger := log.New(multiWriter, ".osz Renamer: ", log.LstdFlags|log.Lshortfile)

	oszDirectory, oszDirectoryReadErr := ioutil.ReadDir(oszDir)

	if oszDirectoryReadErr != nil {
		logger.Fatalf("Failed to read osz directory!")
	}

	for _, file := range oszDirectory {
		if file.IsDir() {
			continue
		}

		if !strings.HasSuffix(file.Name(), ".osz") {
			continue
		}

		oldName := oszDir + "/" + file.Name()
		newName := oszDir + "/" + strings.Split(file.Name(), " ")[0] + ".osz"

		renameErr := os.Rename(oldName, newName)

		if renameErr != nil {
			logger.Printf("Failed to rename %s to %s", oldName, newName)
		}
	}
}
