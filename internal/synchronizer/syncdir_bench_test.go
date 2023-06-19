package synchronizer

import (
	lg "go-base-final/internal/logger"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

type tempPaths struct {
	tempSource string
	tempTarget string
}

var result = func() *tempPaths {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	return &tempPaths{
		tempSource: os.Getenv("BENCHSOURSE"),
		tempTarget: os.Getenv("BENCHTARGET"),
	}
}()

var pathToCopy *CopyPaths = &CopyPaths{
	inPath:  result.tempSource,
	outPath: result.tempTarget,
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
