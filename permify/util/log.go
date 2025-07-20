package util

import (
	"log"
	"os"
)

var (
	LogFile *os.File
	Log     *log.Logger
)

func InitLogger() {
	var err error
	LogFile, err = os.OpenFile("gen.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file %v", err)
	}
	Log = log.New(LogFile, "", log.LstdFlags|log.Lshortfile)
}
