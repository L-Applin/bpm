package utils

import (
	"encoding/base64"
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestBase64Reader(t *testing.T) {
	str := "Hello, Sailor!"
	encodedStr := "SGVsbG8sIFNhaWxvciE="
	reader := strings.NewReader(str)
	b64Reader := Base64reader(reader)
	encodedBytes, err := io.ReadAll(b64Reader)
	if err != nil {
		t.Errorf("error while ReadAll")
		return
	}
	encodedString := string(encodedBytes)
	if encodedString != encodedStr {
		t.Errorf("expected encoded string to be '%s' but was '%s'\n", encodedStr, encodedString)
		return
	}
	fmt.Println(encodedString)
	decodedBytes, err := io.ReadAll(base64.NewDecoder(base64.StdEncoding, strings.NewReader(encodedString)))
	if err != nil {
		t.Errorf("error while ReadAll")
		return
	}
	fmt.Println(string(decodedBytes))
	if string(decodedBytes) != str {
		t.Errorf("expected decoded string to be '%s' but was '%s'\n", str, string(decodedBytes))
		return
	}
}
