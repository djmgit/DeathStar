
// functionpackager contains methods for packaging our lambda function
package functionpackager

import (
	"archive/zip"
	"github.com/rs/zerolog"
	"io"
	"os"
	"path/filepath"
)

// Packager is the function that will be externally aclled to package
// our lambda function and return the path to the created zip file
func Package(deathLogger zerolog.Logger) string {

	// Get the absolute path to he executable
	executablePath, _ := os.Executable()

	// The zip file will be created in the same directory where the executable is present
	dirPath := filepath.Dir(executablePath)

	zipFileName := "lambda_func.zip"

	// Absolute path of the final zip file
	// this path will be returned once the zip file has beem created
	zipFilePath := filepath.Join(dirPath, zipFileName)

	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		deathLogger.Fatal().Msg(err.Error())
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	addToZipFile(zipWriter, executablePath, deathLogger)

	return zipFilePath
}

func addToZipFile(zipWriter *zip.Writer, filePath string, deathLogger zerolog.Logger) {
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

	// Change to deflate to gain better compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
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

