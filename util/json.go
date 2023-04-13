package util

import "encoding/json"

func UnmarshalJson[T any](placeholder *T) func([]byte) (*T, error) {
	return func(data []byte) (*T, error) {
		err := json.Unmarshal(data, placeholder)
		if err != nil {
			return nil, err
		}
		return placeholder, nil
	}
}

func UnmarshalJsonDeref[T any](placeholder *T) func([]byte) (T, error) {
	return func(data []byte) (T, error) {
		var noop T
		ret, err := UnmarshalJson(placeholder)(data)
		if err != nil {
			return noop, err
		}
		return *ret, nil
	}
}

func MarshalJson[T any](data T) ([]byte, error) {
	return json.Marshal(data)
}
