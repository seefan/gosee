package main

import (
	"flag"
	"fmt"
	"io"
	"strings"
	//	"go/ast"
	"github.com/seefan/goerr"
	"go/parser"
	"go/token"
	//	"io/ioutil"
	//	"log"
	"os"
	"path/filepath"
	"sort"
)

var (
	file          = ""
	path          = "./"
	prefix        = "github.com;code.google.com"
	onlyNotExists = false
)

func init() {
	flag.StringVar(&file, "file", file, "请输入要保存的文件名称")
	flag.StringVar(&prefix, "prefix", prefix, "请输入要查询出的包名前缀")
	flag.BoolVar(&onlyNotExists, "notexists", onlyNotExists, "是否只列出不存在的包")
	flag.StringVar(&path, "path", path, "请输入要查询的路径，默认为当前目录。")
}
func main() {
	flag.Parse()
	var f io.Writer
	if len(file) == 0 {
		f = os.Stdout
	} else {
		if ff, err := os.Create(file); err != nil {
			fmt.Println(goerr.NewError(err, "创建输出文件时出错"))
		} else {
			f = ff
			defer ff.Close()
		}
	}
	imports := make(map[string]bool)
	prefixs := strings.Split(prefix, ";")
	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if info == nil || info.IsDir() {
			return nil
		}
		r, err := filepath.Rel(path, p)
		if err != nil {
			return err
		}
		ext := filepath.Ext(r)
		if ext != ".go" {
			return nil
		}
		tms, err := getImport(info.Name(), p)
		if err != nil {
			return err
		}

		for tm, _ := range tms {
			for _, pre := range prefixs {

				if strings.HasPrefix(tm, pre) {
					imports[tm] = false
				}
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println(goerr.NewError(err, "查询导入包时出错"))
	}
	outs := []string{}
	for k, _ := range imports {
		outs = append(outs, k)
	}
	sort.Strings(outs)
	for _, k := range outs {
		_, err := f.Write([]byte(k + "\n"))
		if err != nil {
			fmt.Println(goerr.NewError(err, "输出导入包结果时出错"))
			break
		}
	}
	//TODO 检查包在本地是否存在
}
func getImport(name, path string) (map[string]bool, error) {
	re := make(map[string]bool)
	f, err := os.Open(path)
	if err != nil {
		return nil, goerr.NewError(err, "打开文件 %s 时出错", path)
	}
	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	ff, err := parser.ParseFile(fset, name, f, 0)
	if err != nil {
		return nil, goerr.NewError(err, "读取文件 %s 时出错", path)
	}
	for _, s := range ff.Imports {
		re[strings.Trim(s.Path.Value, `"`)] = false
	}
	return re, nil
}
