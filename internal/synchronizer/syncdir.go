package synchronizer

import (
	"context"
	lg "go-base-final/internal/logger"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// Структура для функции копирования, с файлом которым надо скопировать и куда скопировать
type CopyPaths struct {
	inPath  string
	outPath string
}

func SyncDirectories(ctx context.Context, pathSourse string, pathDestination string, wg *sync.WaitGroup, loggers map[lg.LoggerKey]*log.Logger) error {
	var dirMutex sync.Mutex
	//канал для передачи путей из функции обхода директорий
	filePaths := make(chan *CopyPaths, 10)

	go func() {
		defer close(filePaths)
		for {
			select {
			case <-ctx.Done():
				return
			case pathToCopy := <-filePaths:
				wg.Add(1)
				go func(pathsToCopy *CopyPaths) {
					defer wg.Done()
					select {
					case <-ctx.Done():
						return
					default:
						//Ставим мьютекс для того что бы не создавались одинковые директории несколко раз
						dirMutex.Lock()
						err := CopyFile(pathToCopy, loggers)
						dirMutex.Unlock()
						if err != nil {
							loggers["error"].Println(err)
							return
						}
					}

				}(pathToCopy)

			}

		}
	}()
	//Коллбэк который используется при обходе директорий
	handlePath := func(path string, file fs.DirEntry, err error) error {
		//Здесь проверяем ошибку которую получаем из функции обхода директорий
		if err != nil {
			loggers["error"].Println(err)
			return err
		}
		//Если не является директорий то создаем полный путь для копирования
		if !file.IsDir() {
			finalDestinationPath, err := GetDestinationPath(path, pathSourse, pathDestination, loggers)
			if err != nil {
				loggers["error"].Println(err)
				return err
			}
			//Проверяет что отсутствует такой путь, если ок то отправляем в канал структурой
			if _, err := os.Stat(finalDestinationPath); os.IsNotExist(err) {
				resPath := CopyPaths{
					inPath:  path,
					outPath: finalDestinationPath,
				}
				filePaths <- &resPath
			} else if err != nil {
				loggers["error"].Println(err)
				return err
			}
		}
		return nil
	}
	//Функция для обхода директорий

	if err := filepath.WalkDir(pathSourse, handlePath); err != nil {
		loggers["error"].Println(err)
		return err
	}
	return nil
}

// Создание пути файла который будет спопирован
func GetDestinationPath(filePath, pathSourse, pathDestination string, loggers map[lg.LoggerKey]*log.Logger) (string, error) {
	rel, err := filepath.Rel(pathSourse, filePath)
	if err != nil {
		loggers["error"].Println(err)
		return "", err
	}
	return filepath.Join(pathDestination, rel), nil
}

func CopyFile(pathToCopy *CopyPaths, loggers map[lg.LoggerKey]*log.Logger) error {
	sourceFile, err := os.Open(pathToCopy.inPath)
	if err != nil {
		loggers["error"].Println(err)
		return err
	}
	defer sourceFile.Close()

	dir := filepath.Dir(pathToCopy.outPath)

	if _, err = os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			loggers["error"].Println(err)
			return err
		}
		loggers["info"].Println("Создание директории | ", dir)

	}

	newFile, err := os.Create(pathToCopy.outPath)
	if err != nil {
		loggers["error"].Println(err)
		return err
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, sourceFile)
	if err != nil {
		loggers["error"].Println(err)
		return err
	}
	infoFile, _ := sourceFile.Stat()
	loggers["info"].Println(
		"Копирования файла | ",
		sourceFile.Name(), "|",
		infoFile.Size(),
	)
	return nil

}
