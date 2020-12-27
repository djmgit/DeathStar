package functionpackager

import (
	"archive/zip"
	"github.com/rs/zerolog"
	"os"
	"path/filepath"
)

func Packager(deathLogger zerolog.Logger) {
	executablePath, _ := os.Executable()
	executableName := filepath.Base(executablePath)
	dirPath := filepath.Dir(executablePath)

	zipFileName := "lambda_func.zip"
	zipFilePath := filepath.Join(dirPath, zipFileName)

	zipFile, err := os.Create(zipFileName)
	if err != nil {
		deathLogger.Fatal().Msg(err.Error())
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()
}

