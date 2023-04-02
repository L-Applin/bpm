package utils

import (
	"encoding/base64"
	"io"
)

func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func MapList[T interface{}, V interface{}](list []T, f func(t T) V) []V {
	arr := make([]V, len(list))
	for i, value := range list {
		arr[i] = f(value)
	}
	return arr
}

func Base64reader(reader io.Reader) io.Reader {
	pr, pw := io.Pipe()
	encodedWriter := base64.NewEncoder(base64.StdEncoding, pw)
	go func() {
		defer pw.Close()
		defer encodedWriter.Close()
		io.Copy(encodedWriter, reader)
	}()
	return pr
}
