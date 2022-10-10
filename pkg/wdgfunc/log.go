package wdgfunc

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	mysqlError = "database introspection failed: "
)

type dbErr struct {
	Detail   string `json:"detail"`
	Examples string `json:"examples"`
}

type ErrorLog struct {
	ErrorType string `json:"errorType"` // datasource等错误
	Detail    string `json:"detail"`    // 详细信息
}

// ParseLog 解析日志
func ParseLog() {
	errMap := make(map[string]string, 0)
	fp, err := os.Open("result.log")
	if err != nil {
		fmt.Println(err) //打开文件错误
		return
	}
	buf := bufio.NewScanner(fp)
	contentArr := make([]string, 0)
	for {
		if !buf.Scan() {
			break //文件读完了,退出for
		}
		line := buf.Text() //获取每一行
		if !strings.Contains(line, mysqlError) {
			continue
		}
		errMap[line] = ""
		contentArr = append(contentArr, line)
	}
	//errIndex := 0
	for i, row := range contentArr {
		fmt.Println(i, row)
		if strings.Contains(row, mysqlError) {
			errMap[row] = contentArr[i+6]
		}
	}
	for s, _ := range errMap {
		fmt.Println(s)
	}
}
