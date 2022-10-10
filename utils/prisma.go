package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"
)

var SchemaMap = make(map[string]map[string]string, 0)

func InitQuerySchema(path string) {
	files := GetDicAllChildFilePath(path)
	dirs := GetDicAllChildDicPath(path)

	for _, dirName := range dirs {
		if dirName == path {
			continue
		}
		// 获取目录名，也就是module名
		dicNameArr := strings.Split(dirName, "/")
		moduleName := dicNameArr[len(dicNameArr)-1]
		moduleMap := make(map[string]string, 0)
		for _, fileName := range files {
			content, _ := ioutil.ReadFile(fileName)
			key := filepath.Base(fileName)
			moduleMap[key] = string(content)
		}
		SchemaMap[moduleName] = moduleMap
	}
}

func GetQuerySchema(module, methodName string, param map[string]interface{}) (result string) {
	querySchema := SchemaMap[module][methodName]
	if querySchema == "" {
		return
	}
	// 遍历参数列表,将占位符替换
	for key, val := range param {
		paramKey := fmt.Sprintf("$%s", key)
		valType := reflect.TypeOf(val)
		if valType.Kind() == reflect.String {
			querySchema = strings.ReplaceAll(querySchema, paramKey, fmt.Sprintf(`"%v"`, val))
			continue
		}
		// 如果是数组
		if valType.Kind() == reflect.Slice {
			queryVal := ""
			// 序列化数组
			arrBytes, _ := json.Marshal(val)
			queryVal = string(arrBytes)
			querySchema = strings.ReplaceAll(querySchema, paramKey, fmt.Sprintf(`"%v"`, queryVal))
			// 将合成后的"[  ]" 转为[]
			querySchema = strings.ReplaceAll(querySchema, "\"[", "[")
			querySchema = strings.ReplaceAll(querySchema, "]\"", "]")
		}
		querySchema = strings.ReplaceAll(querySchema, paramKey, fmt.Sprintf(`%v`, val))
	}

	//result = fmt.Sprintf("{\"query\": \"%s\",\"variables\": {} }", querySchema)
	result = querySchema
	return
}

// GetDefaultDBSchema 默认数据源的schema
func GetDefaultDBSchema() (result string) {
	dbPath := filepath.Join(GetRootPath(), GetDefaultOauthDbPath())
	outputPath := filepath.Join(GetRootPath(), GetPrismaDBPath())
	return fmt.Sprintf(`
datasource db {
	// could be postgresql or mysql
	provider = "sqlite"
	url      = "file:%s"
}
generator db {
	provider = "go run github.com/prisma/prisma-client-go"
	// set the output folder and package name
	output           = "%s"
	package          = "prismaDB"
	disableGoBinaries = "true"
}
`, dbPath, outputPath)
}
