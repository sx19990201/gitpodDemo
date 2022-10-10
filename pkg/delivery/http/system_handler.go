package http

import (
	"bufio"
	"fmt"
	"github.com/fire_boom/domain"
	"github.com/fire_boom/utils"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"sync"
	"time"
)

const (
	Start = iota + 1
	ReStart
	Close
)
const (
	StartSuccess = iota + 1
	StartFail
	CloseSuccess
	CloseFail
	StartKey = "-start"
	CloseKey = "-close"
)

var Signal = make(chan<- int64, 1)

var StartSignal = make(chan<- int64, 1)
var CloseSignal = make(chan<- int64, 1)

type WdgSystemHandler struct {
}

func InitWdgSystemHandler(e *echo.Echo) {
	NewScriptHandler(e)
}

func NewWdgSystemHandler(e *echo.Echo) {
	//handler := &WdgSystemHandler{}
	//v1 := e.Group("/api/v1")
	//{
	//	system := v1.Group("/wdg")
	//	{
	//		//system.GET("/start", handler.Start)
	//		//system.GET("/close", handler.Close)
	//		//system.GET("/reStart", handler.ReStart)
	//	}
	//}
}

// Start @Title Start
// @Description 启动
// @Accept  json
// @Tags  wdg
// @Success 200 {object} domain.SystemConfig "启动成功"
// @Failure 400	"启动失败"
// @Router /api/v1/wdg/start  [GET]
func (s *WdgSystemHandler) Start(c echo.Context) (err error) {
	setting, err := utils.GetSettingConfig()
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.SettingCode, err))
	}
	return c.JSON(http.StatusOK, SuccessWriteResult(setting.System))
}

// Close @Title Close
// @Description 关闭
// @Accept  json
// @Tags  wdg
// @Success 200 {object} domain.SystemConfig "关闭成功"
// @Failure 400	"关闭失败"
// @Router /api/v1/wdg/close  [GET]
func (s *WdgSystemHandler) Close(c echo.Context) (err error) {
	setting, err := utils.GetSettingConfig()
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.SettingCode, err))
	}
	return c.JSON(http.StatusOK, SuccessWriteResult(setting.System))
}

// ReStart @Title ReStart
// @Description 重启
// @Accept  json
// @Tags  wdg
// @Success 200 {object} domain.SystemConfig "重启成功"
// @Failure 400	"重启失败"
// @Router /api/v1/wdg/reStart  [GET]
func (s *WdgSystemHandler) ReStart(c echo.Context) (err error) {
	setting, err := utils.GetSettingConfig()
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.SettingCode, err))
	}
	return c.JSON(http.StatusOK, SuccessWriteResult(setting.System))
}

var execMap sync.Map

func GetExecCommand(key string) *exec.Cmd {
	if val, ok := execMap.Load(key); ok {
		// 存在直接返回
		cmd, ok := val.(*exec.Cmd)
		if !ok {
			return nil
		}
		return cmd
	}
	return nil
}

func StoreExecCommand(key, cmdStr string) *exec.Cmd {
	cmdType := "bash"
	if runtime.GOOS == "windows" {
		cmdType = "cmd"
	}
	cmd := exec.Command(cmdType, key, fmt.Sprintf("cd wundergraph/.wundergraph && %s", cmdStr))

	execMap.Store(key, cmd)
	return cmd
}
func CheckStartState(state *os.ProcessState, startSignal chan int64) {
	// 判断是否启动成功
	for {
		// 为空说明进程没有运行了
		if state == nil {
			startSignal <- StartFail
			break
		}
		// 扫描9991端口，一秒超时时间
		if ScanPort("127.0.0.1", 9991) {
			// 9991启动成功
			startSignal <- StartSuccess
			break
		}
	}
}

// ScanPort 扫描端口通过ip:端口检测端口是否被占用
func ScanPort(hostname string, port int) bool {
	fmt.Printf("scanning port %d \n", port)
	p := strconv.Itoa(port)
	addr := net.JoinHostPort(hostname, p)
	conn, err := net.DialTimeout("tcp", addr, 1*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func InitWundergraph(v *viper.Viper) {
	// windows 系统可能要把 bash 改为 cmd?
	cmdType := "bash"
	if runtime.GOOS == "windows" {
		cmdType = "cmd"
	}

	c := exec.Command(cmdType, "-c", "cd wundergraph/.wundergraph && wunderctl up --debug --listen-addr 0.0.0.0:9991")
	//c := exec.Command(cmdType, "-c", " cd wundergraph/.wundergraph && wunderctl up --debug --listen-addr 0.0.0.0:9991")
	stdout, err := c.StdoutPipe()
	if err != nil {
		panic(fmt.Sprintf("init wunderGraph stdoutPipe meet err=%+v", err))
	}

	file, err := os.OpenFile(utils.GetWdgLogPath(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("Failed to log to file, using default stderr")
	}

	go func() {
		reader := bufio.NewReader(stdout)
		for {
			readString, err := reader.ReadBytes('\n')
			if err != nil || err == io.EOF {
				break
			}
			file.Write(readString)
		}
	}()
	if err := c.Run(); err != nil {
		panic(fmt.Sprintf("init wunderGraph run cmd meet err=%+v", err))
	}

}
