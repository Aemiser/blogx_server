package hash

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

func Md5(data []byte) string {
	md5New := md5.New()
	md5New.Write(data)
	return hex.EncodeToString(md5New.Sum(nil))
}

func FileMd5(filePath string) (h string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	byteDate, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	h = Md5(byteDate)
	return
}
