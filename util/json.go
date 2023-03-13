package util

import "encoding/json"

func UnmarshalJson[T any](placeholder T) func([]byte) (T, error) {
	return func(data []byte) (T, error) {
		var noop T
		err := json.Unmarshal(data, placeholder)
		if err != nil {
			return noop, err
		}
		return placeholder, nil
	}
}

func UnmarshalJsonStruct[T any](data []byte) (T, error) {
	var placeholder T
	return UnmarshalJson(placeholder)(data)
}

func UnmarshalJsonStructPtr[T any](data []byte) (*T, error) {
	var placeholder T
	ret, err := UnmarshalJson(placeholder)(data)
	return &ret, err
}
