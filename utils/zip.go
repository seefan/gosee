package utils

import (
	"archive/zip"
	"strings"
	//	"bytes"
	"github.com/seefan/goerr"
	"io"
	"os"
	"path/filepath"
)

func UnZipDir(file, dir string) error {
	// Open a zip archive for reading.
	r, err := zip.OpenReader(file)
	if err != nil {
		return goerr.NewError(err, "打开压缩包%s时出错", file)
	}
	defer r.Close()
	psr := string(os.PathSeparator)
	if !strings.HasSuffix(dir, psr) {
		dir += psr
	}

	// Iterate through the files in the archive,
	// printing some of their contents.
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return goerr.NewError(err, "解压%s时出错", file)
		}
		defer rc.Close()
		MakeDirs(dir + f.Name)
		dstFile, err := os.Create(dir + f.Name)
		if err != nil {
			return goerr.NewError(err, "解压%s时出错", file)
		}
		_, err = io.Copy(dstFile, rc)
		if err != nil {
			return goerr.NewError(err, "解压%s时出错", file)
		}
	}
	return nil
}
func ZipDir(dir, file string) error {
	// Create a buffer to write our archive to.
	tmpFile := os.TempDir() + file
	err := MakeDirs(tmpFile)
	if err != nil {
		return goerr.NewError(err, "创建压缩包%s时出错", file)
	}
	dstFile, err := os.Create(tmpFile)
	if err != nil {
		return goerr.NewError(err, "创建压缩包%s时出错", file)
	}
	defer dstFile.Close()
	// Create a new zip archive.
	w := zip.NewWriter(dstFile)
	//psr := string(os.PathSeparator)
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info == nil || info.IsDir() {
			return nil
		}
		f, err := w.Create(path)
		if err != nil {
			return goerr.NewError(err, "创建压缩文件%s时出错", path)
		}
		srcFile, err := os.Open(path)
		if err != nil {
			return goerr.NewError(err, "打开待压缩文件%s时出错", path)
		}
		defer srcFile.Close()
		_, err = io.Copy(f, srcFile)
		if err != nil {
			return goerr.NewError(err, "压缩文件时出错")
		}
		return nil
	})
	if err != nil {
		defer w.Close()
		return err
	}

	// Make sure to check the error on Close.
	err = w.Close()
	if err != nil {
		return goerr.NewError(err, "保存文件时出错")
	}
	MakeDirs(file)
	err = CopyFile(tmpFile, file)
	if err != nil {
		return goerr.NewError(err, "复制临时文件到%s时出错", file)
	}
	os.Remove(tmpFile)
	return nil
}
