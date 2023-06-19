package synchronizer

import (
	lg "go-base-final/internal/logger"
	"log"
	"os"
	"testing"
)

// Пути к файлам захардкожены для тестирования производительности
var pathToCopy *CopyPaths = &CopyPaths{
	inPath:  `D:\golang\GO-BASE-FINAL\sourseDirectiry\others\futerman\Future.Man.S03E01.1080p.AMZN.WEB-DL.H.264.RGzsRutracker.[Wentworth_Miller].NTb.mkv`,
	outPath: `D:\golang\GO-BASE-FINAL\targetDirectory\others\futerman\Future.Man.S03E01.1080p.AMZN.WEB-DL.H.264.RGzsRutracker.[Wentworth_Miller].NTb.mkv`,
}

type benchStruct struct {
	loggers map[lg.LoggerKey]*log.Logger
	file    *os.File
}

func initLoggers() *benchStruct {
	infoLogger, errLoger, err := lg.InitLoggers()
	if err != nil {
		log.Println(err)
	}

	infoFile := infoLogger.Writer().(*os.File)

	loggers := map[lg.LoggerKey]*log.Logger{
		lg.LoggerKey("info"):  infoLogger,
		lg.LoggerKey("error"): errLoger,
	}

	benchData := &benchStruct{
		loggers: loggers,
		file:    infoFile,
	}

	return benchData

}

var loggers *benchStruct = initLoggers()

func BenchmarkCopiFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CopyFile(pathToCopy, loggers.loggers)
		defer func() {
			if err := loggers.file.Close(); err != nil {
				log.Fatal(err)
			}
		}()
	}
}
