package json_util

import "encoding/json"

func Unmarshal[T any, S ~[]byte | ~string](data S) (v T, err error) {
	err = json.Unmarshal([]byte(data), &v)
	return
}

func MarshalString(v any) (string, error) {
	b, err := json.Marshal(v)
	return string(b), err
}
