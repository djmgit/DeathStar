package functionpackager

import (
	"os"
	"path/filepath"
)

func Packager() {
	executablePath, _ := os.Executable()
	executableName := filepath.Base(executablePath)
	dirPath := filepath.Dir(executablePath)

	zipFileName := "lambda_func.zip"
}

