package helpers

import (
	"io"
	"log"
	"os"
)

func SetupLogger() (io.Writer, func () error, error) {
	logFile, err := os.OpenFile("log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatal(err)
		return nil, nil, err
	}

	logMW := io.MultiWriter(logFile, os.Stdout)
	log.SetOutput(logMW)
	log.SetFlags(log.LstdFlags|log.Lmsgprefix)
	
	return logMW, logFile.Close, nil
}
