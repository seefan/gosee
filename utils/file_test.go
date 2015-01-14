package utils

import (
	"testing"
)

func Test_exists(t *testing.T) {
	ok := ExistsFile("file.go")
	t.Log(ok)
	ok = ExistsFile("file12.go")
	t.Log(ok)
}

func Test_copyfile(t *testing.T) {
	err := CopyFile("file.go", "file1.go")
	if err != nil {
		t.Error(err)
	}
}
