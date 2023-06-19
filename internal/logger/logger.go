package logger

import (
	"log"
	"os"
)

type LoggerKey string

// Иннициализация логгеров
func InitLoggers() (infoLog, errorLog *log.Logger, err error) {
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, nil, err
	}

	//Разные логеры для событий
	infoLog = log.New(f, "INFO\t", log.Ldate|log.Ltime)
	errorLog = log.New(f, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	return infoLog, errorLog, nil
}
