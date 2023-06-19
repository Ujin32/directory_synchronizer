package main

import (
	"context"
	"flag"
	lg "go-base-final/internal/logger"
	sh "go-base-final/internal/synchronizer"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	pathSourse      string
	pathDestination string
	syncPeriod      string
)

func main() {
	//Флаги командной строки
	flag.StringVar(&pathSourse, "output", `\golang\GO-BASE-FINAL\sourseDirectiry\`, "Intput file path")
	flag.StringVar(&pathDestination, "input", `D:\golang\GO-BASE-FINAL\targetDirectory\`, "Output file path")
	flag.StringVar(&syncPeriod, "syncPeriod", `10`, "Directory synchronization period")
	flag.Parse()

	//Контекст на отмену
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//Инициализация логгеров
	infoLogger, errLoger, err := lg.InitLoggers()
	if err != nil {
		log.Fatal(err)
	}

	//Получение файлового дескриптора из логгера
	infoFile := infoLogger.Writer().(*os.File)

	//Закрытие файлового дескриптора
	defer func() {
		if err := infoFile.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	//Создаю мапу с логерами для передачи этих двух логов в функции
	loggers := map[lg.LoggerKey]*log.Logger{
		lg.LoggerKey("info"):  infoLogger,
		lg.LoggerKey("error"): errLoger,
	}

	syncPeriod, err := strconv.Atoi(syncPeriod)
	if err != nil {
		log.Fatal(err)
	}

	//Примитив для синхронизации горутин которые копируют файлы
	var wg sync.WaitGroup

	//Запуск функции синхронизации
	for {
		if err := sh.SyncDirectories(ctx, pathSourse, pathDestination, &wg, loggers); err != nil {
			loggers["error"].Println("Синхронизация не удалась ", err)

		}
		time.Sleep(time.Second * time.Duration(syncPeriod))
		//Ожидание всех горутин которые копируют файлы
		wg.Wait()
	}

}
