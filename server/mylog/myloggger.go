package mylog

import (
	"log"
	"os"
)

var (
	Debug *log.Logger
	Info  *log.Logger
	Error *log.Logger
)

func InitLogger() {
	log.Println("Init MyLogger...")

	Debug = log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stdout, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
}
