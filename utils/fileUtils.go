package utils

import (
	"bufio"
	"github.com/nasa9084/go-openapi"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func CheckOASFileContent(path string) (err error) {
	doc, err := openapi.LoadFile(path)
	if err != nil {
		return
	}
	return doc.Validate()
}

func GetDicAllChildFilePath(path string) (result []string) {
	//获取当前目录下的所有文件或目录信息
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			result = append(result, filepath.ToSlash(path))
		}
		return nil
	})
	return
}

// GetDicAllChildDicPath 获取所有文件夹，路径平铺展开
func GetDicAllChildDicPath(path string) (result []string) {
	//获取当前目录下的所有文件或目录信息
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			result = append(result, filepath.ToSlash(path))
		}
		return nil
	})
	return
}

func ReadFileLine(path string) (result []string) {
	// 如果文件不存在
	if !FileExist(path) {
		// 创建文件并写入空
		WriteFile(path, "")
	}
	f, _ := os.Open(path)
	defer f.Close()

	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行

		if err != nil || io.EOF == err {
			if line == "" {
				break
			}
		}
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}
		result = append(result, line)
	}
	return
}

func WriteFile(path, content string) error {
	dicPath := filepath.Dir(path)
	_, err := os.Stat(dicPath)
	// 文件不存在
	if os.IsNotExist(err) {
		// 创建
		err = os.MkdirAll(dicPath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	err = ioutil.WriteFile(path, []byte(content), 0644)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := recover(); err != nil {
			log.Error("file write fail ,err : ", err)
			return
		}
	}()
	return nil
}

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}
