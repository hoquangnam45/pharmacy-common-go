package common

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"github.com/hoquangnam45/pharmacy-common-go/helper/errorHandler"
)

func FileReadAllBytes(file *os.File) *errorHandler.MaybeError[[]byte] {
	return errorHandler.Transform(func() ([]byte, error) {
		return io.ReadAll(file)
	})
}

func CloseFile(file *os.File) {
	file.Close()
}

func UnmarshalBytesToMap(byteValue []byte) *errorHandler.MaybeError[map[string]any] {
	return errorHandler.Transform(func() (map[string]any, error) {
		var result map[string]any
		err := json.Unmarshal([]byte(byteValue), &result)
		return result, err
	})
}

func FatalLog(err error) {
	log.Fatal(err)
}
