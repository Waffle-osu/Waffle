package helpers

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

var Logger *log.Logger
var Guard *log.Logger

func InitializeLogger() {
	filename := fmt.Sprintf("logs/%d-log.txt", time.Now().Unix())
	filenameGuard := fmt.Sprintf("logs/%d-log-guard.txt", time.Now().Unix())

	file, fileErr := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	fileGuard, fileErrGuard := os.OpenFile(filenameGuard, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)

	if fileErr != nil {
		panic(fileErr)
	}

	if fileErrGuard != nil {
		panic(fileErr)
	}

	multiWriter := io.MultiWriter(file, os.Stdout)
	multiWriterGuard := io.MultiWriter(fileGuard, os.Stdout)

	Logger = log.New(multiWriter, "Waffle: ", log.LstdFlags|log.Lshortfile)
	Guard = log.New(multiWriterGuard, "WaffleGuard: ", log.LstdFlags)
}
