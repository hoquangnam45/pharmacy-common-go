package util

import (
	"encoding/json"
	"io"
	"os"
)

func FileReadAllBytes(file *os.File) ([]byte, error) {
	return io.ReadAll(file)
}

func CloseFile(file *os.File) {
	file.Close()
}

func UnmarshalBytesToMap(byteValue []byte) (map[string]any, error) {
	var result map[string]any
	err := json.Unmarshal([]byte(byteValue), &result)
	return result, err
}

func FatalLog(err error) {
	Logger.Panic(err.Error())
}

func Panic(err error) {
	panic(err)
}

func ReadAllThenClose(r io.ReadCloser) ([]byte, error) {
	defer r.Close()
	return io.ReadAll(r)
}
