package utils

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"github.com/fire_boom/domain"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"github.com/ghodss/yaml"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

const (
	DateFormat     = "2006-01-02"
	DateTimeFormat = "2006-01-02 15:04:05"
)

func WriteCVS(path string, data [][]string) error {
	rowContent := ""
	for _, row := range data {
		//lineContent := ""
		//for _, line := range row {
		//	tempLine := ""
		//	if strings.Contains(line, "\"") {
		//		tempLine = strings.ReplaceAll(line, `"`, `\"`)
		//	}
		//	tempLine = fmt.Sprintf(`"%s"`, tempLine)
		//	strings.
		//	lineContent = fmt.Sprintf("%s %s,", lineContent, tempLine)
		//}
		line := strings.Join(row, ",")

		rowContent = fmt.Sprintf("%s %s \n", rowContent, line)
	}
	err := WriteFile(path, rowContent)
	if err != nil {
		return err
	}
	return nil
}

func StrContainsSpace(val string) bool {
	return strings.Contains(val, " ")
}

func GetSettingConfig() (result domain.Setting, err error) {
	path := GetSettingPath()
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Error("GetSettingConfig ioutil.ReadFile err : ", err)
		err = domain.FileReadErr
		return
	}
	if len(content) == 0 {
		return
	}
	err = json.Unmarshal(content, &result)
	if err != nil {
		log.Error("GetSettingConfig json.Unmarshal err : ", err)
		err = domain.JsonUnMarshalErr
		return
	}
	return
}

func GetEnvVal(envs []domain.FbEnv) (map[string]string, error) {
	result := make(map[string]string, 0)
	// TODO 判断当前环境，取当前环境的值
	setting, err := GetSettingConfig()
	if err != nil {
		log.Error("failed to get current environment")
		return result, err
	}
	if len(envs) > 0 {
		for _, env := range envs {
			// 开发环境
			if setting.System.DevSwitch == true {
				result[env.Key] = env.DevEnv
				continue
			}
			result[env.Key] = env.ProEnv
		}
	}
	return result, nil
}

func ZipFiles(filename string, files []string) error {

	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {
		if err = AddFileToZip(zipWriter, file); err != nil {
			return err
		}
	}
	return nil
}

func AddFileToZip(zipWriter *zip.Writer, filename string) error {

	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// Using FileInfoHeader() above only uses the basename of the file. If we want
	// to preserve the folder structure we can overwrite this with the full path.
	filename = strings.ReplaceAll(filename, GetSdkSrcPath(), "")
	header.Name = filename

	// Change to deflate to gain better compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}

func GetRootPath() string {
	dir := GetCurrentAbPathByExecutable()
	tmpDir, _ := filepath.EvalSymlinks(os.TempDir())
	if strings.Contains(dir, tmpDir) {
		return GetCurrentAbPathByCaller()
	}
	return dir
}

func GetIndexDBID() int64 {
	content := struct {
		ID int64 `json:"id"`
	}{}
	bytes, _ := ioutil.ReadFile(GetIndexDBIDJsonPath())
	json.Unmarshal(bytes, &content)
	if content.ID == 0 {
		content.ID = -1
	}
	return content.ID
}

func SetIndexDBID(id int64) {
	content := struct {
		ID int64 `json:"id"`
	}{ID: id}
	bytes, _ := json.Marshal(content)
	WriteFile(GetIndexDBIDJsonPath(), string(bytes))
}

// GetCurrentAbPathByExecutable 获取当前执行文件绝对路径
func GetCurrentAbPathByExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}

// GetCurrentAbPathByCaller 获取当前执行文件绝对路径（go run）
func GetCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}

func CheckJsonStr(jsonStr []byte) bool {
	return json.Valid(jsonStr)
}

func CheckYamlStr(yamlStr []byte) bool {
	if _, err := yaml.YAMLToJSON(yamlStr); err != nil {
		return false
	}
	return true
}

func GenerateUUID() string {
	return uuid.NewV4().String()
}

func ArrStrFormatStrArr(arrStr string) (result []string, err error) {
	err = json.Unmarshal([]byte(arrStr), &result)
	return
}

func ArrStrFormatIntArr(arrStr string) (result []int64, err error) {
	err = json.Unmarshal([]byte(arrStr), &result)
	return
}

func ArrStrFormatFloatArr(arrStr string) (result []float64, err error) {
	err = json.Unmarshal([]byte(arrStr), &result)
	return
}

func Empty(val interface{}) bool {
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.String, reflect.Array:
		return v.Len() == 0
	case reflect.Map, reflect.Slice, reflect.Chan:
		return v.Len() == 0 || v.IsNil()
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Complex64, reflect.Complex128:
		return v.Complex() == 0+0i
	case reflect.Func:
		return v.IsNil()
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	case reflect.Struct:
		return reflect.DeepEqual(val, reflect.Zero(v.Type()).Interface())
	case reflect.UnsafePointer:
		return reflect.DeepEqual(val, reflect.Zero(v.Type()).Interface())
	case reflect.Invalid:
		return false
	}
	return reflect.DeepEqual(val, reflect.Zero(v.Type()).Interface())
}

func InterfaceToString(i interface{}) string {
	switch v := i.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case bool:
		return strconv.FormatBool(v)
	default:
		bytes, _ := json.Marshal(v)
		return string(bytes)
	}
}
