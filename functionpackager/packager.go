package functionpackager

import (
	"archive/zip"
	"github.com/rs/zerolog"
	"io"
	"os"
	"path/filepath"
)

func Packager(deathLogger zerolog.Logger) {
	executablePath, _ := os.Executable()
	dirPath := filepath.Dir(executablePath)

	zipFileName := "lambda_func.zip"
	zipFilePath := filepath.Join(dirPath, zipFileName)

	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		deathLogger.Fatal().Msg(err.Error())
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	AddToZipFile(zipWriter, executablePath, deathLogger)
}

func AddToZipFile(zipWriter *zip.Writer, filePath string, deathLogger zerolog.Logger) {
	fileName := filepath.Base(filePath)

	fileToZip, err := os.Open(filePath)
	if err != nil {
		deathLogger.Fatal().Msg(err.Error())
	}

	info, err := fileToZip.Stat()
	if err != nil {
		deathLogger.Fatal().Msg(err.Error())
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		deathLogger.Fatal().Msg(err.Error())
	}

	header.Name = fileName

	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		deathLogger.Fatal().Msg(err.Error())
	}
	_, err = io.Copy(writer, fileToZip)
	if err != nil {
		deathLogger.Fatal().Msg(err.Error())
	}
}

