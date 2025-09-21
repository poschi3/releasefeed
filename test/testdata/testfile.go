package testdata

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func GetFilePath(elem ...string) string {
	elem = append([]string{basepath}, elem...)
	return filepath.Join(elem...)
}

func GetFileContent(elem ...string) ([]byte, error) {

	absolutPath := GetFilePath(elem...)

	file, err := os.Open(absolutPath)
	if err != nil {
		return nil, err

	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return data, nil

}
