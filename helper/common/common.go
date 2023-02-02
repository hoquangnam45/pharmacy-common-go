package common

import (
	"encoding/json"
	"io"
	"log"
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
	log.Fatal(err)
}
