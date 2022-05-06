package helpers

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

var Logger *log.Logger

func InitializeLogger() {
	filename := fmt.Sprintf("logs/%d-log.txt", time.Now().Unix())

	file, fileErr := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if fileErr != nil {
		panic(fileErr)
	}

	multiWriter := io.MultiWriter(file, os.Stdout)

	Logger = log.New(multiWriter, "Waffle: ", log.LstdFlags|log.Lshortfile)
}
