package http

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// VerOption 全局公共变量：选项
var VerOption Option

// Option 选项
type Option struct {
	Port    string `yaml:"port"`    // 端口
	ApiHost string `yaml:"apiHost"` // API主机（IP + Port）
}

// DoProxy 转发请求
func DoProxy(r *http.Request, host string) (result string, err error) {
	// 创建一个HttpClient用于转发请求
	cli := &http.Client{}

	// 读取请求的Body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("读取请求体发生错误")
		// 响应状态码
		return
	}

	// 转发的URL
	//reqURL := option.ApiHost + r.URL.String()
	fmt.Println()
	reqURL := fmt.Sprintf("%s%s", host, r.RequestURI)
	reqURL = strings.ReplaceAll(reqURL, ":9123", ":9991")

	// 创建转发用的请求
	reqProxy, err := http.NewRequest(r.Method, reqURL, strings.NewReader(string(body)))
	if err != nil {
		log.Println("创建转发请求发生错误")
		// 响应状态码
		return
	}

	// 转发请求的 Header
	for k, v := range r.Header {
		reqProxy.Header.Set(k, v[0])
	}

	// 发起请求
	responseProxy, err := cli.Do(reqProxy)
	if err != nil {
		log.Println("转发请求发生错误")
		// 响应状态码
		return
	}
	defer responseProxy.Body.Close()

	// 转发响应的Body数据
	var data []byte

	// 读取转发响应的Body
	data, err = ioutil.ReadAll(responseProxy.Body)
	if err != nil {
		log.Println("读取响应体发生错误")
		// 响应状态码
		return
	}

	// 转发响应的输出数据
	var dataOutput []byte
	// gzip压缩编码数据
	dataOutput = data
	// 打印转发响应的Body数据，查看转发响应的响应数据时需要。
	result = string(dataOutput)

	// response的Body不能多次读取，
	// 上面已经被读取过一次，需要重新生成可读取的Body数据。
	//resProxyBody := ioutil.NopCloser(bytes.NewBuffer(data))
	//defer resProxyBody.Close() // 延时关闭
	//
	//// 响应状态码
	//w.WriteHeader(responseProxy.StatusCode)
	//// 复制转发的响应Body到响应Body
	//io.Copy(w, resProxyBody)
	return
}
