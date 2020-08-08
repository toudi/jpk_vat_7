package common

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

func FileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func Sha256File(filename string) []byte {
	file, _ := os.Open(filename)
	hash := sha256.New()

	io.Copy(hash, file)

	return hash.Sum(nil)
}

func Md5File(filename string) []byte {
	file, _ := os.Open(filename)
	hash := md5.New()

	io.Copy(hash, file)

	return hash.Sum(nil)
}

func FileSize(filename string) (int64, error) {
	fileStat, err := os.Stat(filename)
	if err != nil {
		return -1, fmt.Errorf("Nie udało się pobrać informacji o pliku: %v", err)
	}

	return fileStat.Size(), nil
}
