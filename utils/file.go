package utils

import (
	"github.com/seefan/goerr"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//拷贝目录
//srcDir 源目录
//dstDir 目标目录
func CopyDir(srcDir, dstDir string) (err error) {
	//首先判断输出目录是否存在该目录，不存在就创建
	if _, err := os.Stat(dstDir); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(dstDir, 0755) //所有人读写
		}
	}
	psr := string(os.PathSeparator)
	if !strings.HasSuffix(srcDir, psr) {
		srcDir += psr
	}
	if !strings.HasSuffix(dstDir, psr) {
		dstDir += psr
	}
	//查看目录下文件，都拷贝过去
	//log.Debugf("copy dir %s to %s", srcDir, dstDir)
	if fis, err := ioutil.ReadDir(srcDir); err == nil {
		for _, fi := range fis {
			if fi.IsDir() { //如果是目录，就也目录下所有文件都拷贝过去
				if err := CopyDir(srcDir+fi.Name(), dstDir+fi.Name()); err != nil {
					return err
				}
			} else {
				if err := CopyFile(srcDir+fi.Name(), dstDir+fi.Name()); err != nil {
					return err
				}
			}
		}
	} else {
		return err
	}
	return nil
}

//文件或目录是否存在
//path 要判断的目标
func ExistsFile(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) { //已经存在的文件，不在重复处理
			return false
		}
	}
	return true
}

//拷贝一个文件
//src 源文件
//dst 目标文件
func CopyFile(src, dstPath string) error {
	srcPath, err := filepath.EvalSymlinks(src)
	if err != nil {
		return goerr.NewError(err, "解析来源文件名%s出错", src)
	}

	srcFile, err := os.Open(srcPath)
	if err != nil {
		return goerr.NewError(err, "打开来源文件%s出错", srcPath)
	}
	defer srcFile.Close()
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return goerr.NewError(err, "创建目标文件%s出错", dstPath)
	}
	defer dstFile.Close()
	if _, err = io.Copy(dstFile, srcFile); err != nil {
		return goerr.NewError(err, "拷贝文件%s到%s出错", srcPath, dstPath)
	}
	return nil
}

//生成不存在的路径
func MakeDirs(file string) error {
	psr := string(os.PathSeparator)
	n := strings.LastIndex(file, psr)
	if n > 0 {
		file = SubString(file, 0, n)
	}
	if ExistsFile(file) {
		return nil
	}
	err := os.MkdirAll(file, 0766)
	if err != nil {
		return goerr.NewError(err, "创建目录%s 时出错", file)
	}
	return nil
}
