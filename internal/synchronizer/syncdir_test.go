package synchronizer

import (
	"fmt"
	lg "go-base-final/internal/logger"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

type MockLogger struct{}

func (l *MockLogger) Write(p []byte) (n int, err error) {
	// Реализация метода Write для интерфейса io.Writer
	return len(p), nil
}

func TestGetDestinationPath(t *testing.T) {
	req := require.New(t)
	loggers := map[lg.LoggerKey]*log.Logger{
		lg.LoggerKey("error"): log.New(new(MockLogger), "", 0),
	}
	sourseDir := createTempDir(t, "sourse")
	targetDir := createTempDir(t, "target")
	sourseFile := creatTempFile(t, sourseDir, "tempFile.txt")

	//Валидный путь
	t.Run("valid path", func(t *testing.T) {
		finalDestinationPath, err := GetDestinationPath(sourseFile, sourseDir, targetDir, loggers)
		require.NoError(t, err)

		pathJoin := filepath.Join(targetDir, "tempFile.txt")
		req.Equal(finalDestinationPath, pathJoin)

	})
	//Не валидный путь
	t.Run("invalid path", func(t *testing.T) {
		invalidFilePath := "nonexistent.txt"
		finalDestinationPath, err := GetDestinationPath(invalidFilePath, sourseDir, targetDir, loggers)
		req.Error(err)
		req.Empty(finalDestinationPath)

	})

}

func createTempDir(t *testing.T, prefix string) string {
	t.Helper()

	tempSourseDir, err := os.MkdirTemp("", prefix)
	require.NoError(t, err)

	return tempSourseDir
}

func creatTempFile(t *testing.T, dir, nameFile string) string {
	t.Helper()

	filePath := filepath.Join(dir, nameFile)
	file, err := os.Create(filePath)
	require.NoError(t, err)
	file.Close()
	return filePath
}

func TestCopyFile(t *testing.T) {
	loggers := map[lg.LoggerKey]*log.Logger{
		lg.LoggerKey("error"): log.New(new(MockLogger), "", 0),
		lg.LoggerKey("info"):  log.New(new(MockLogger), "", 0),
	}
	//Успешное копирование файла
	t.Run("Sucsesfull CopyFile", func(t *testing.T) {
		req := require.New(t)

		fmt.Println(req, loggers)
		sourseDir := createTempDir(t, "sourse")
		targetDir := createTempDir(t, "target")

		sourseFilePath := creatTempFile(t, sourseDir, "tempFile.txt")
		fileBase := filepath.Base(sourseFilePath)
		destinationFilePath := filepath.Join(targetDir, fileBase)
		resPath := &CopyPaths{
			inPath:  sourseFilePath,
			outPath: destinationFilePath,
		}

		err := CopyFile(resPath, loggers)
		req.NoError(err)

		sourseFileData, _ := os.ReadFile(sourseFilePath)
		targetFileData, _ := os.ReadFile(destinationFilePath)

		req.Equal(sourseFileData, targetFileData)

	})
	//Не правильное имя файла
	t.Run("Failed: 'targetDir' to CopyFile", func(t *testing.T) {
		req := require.New(t)
		sourseDir := createTempDir(t, "sourse")
		targetDir := createTempDir(t, "target")
		sourseFilePath := creatTempFile(t, sourseDir, "tempFile.txt")
		destinationFilePath := filepath.Join(targetDir, "")
		resPath := &CopyPaths{
			inPath:  sourseFilePath,
			outPath: destinationFilePath,
		}

		err := CopyFile(resPath, loggers)
		req.Error(err)
	})
}
