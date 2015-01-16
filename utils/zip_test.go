package utils

import (
	"testing"
)

func Test_zip(t *testing.T) {
	if err := ZipDir("./", "./1.zip"); err != nil {
		t.Error(err)
	}
}
func Test_unzip(t *testing.T) {
	if err := UnZipDir("./1.zip", "./abc"); err != nil {
		t.Error(err)
	}
}
