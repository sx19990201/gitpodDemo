package wdgfunc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/fire_boom/domain"
	"github.com/fire_boom/utils"
	"github.com/wundergraph/wundergraph/cli/commands"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type wunderClient struct {
	WunderCtlClient commands.API    `json:"wunderCtlClient"`
	Context         Context         `json:"context"`
	RootContext     context.Context `json:"rootContext"`
}

type Context struct {
	Ctx    context.Context
	Cancel context.CancelFunc
}

var wunderCtlClient wunderClient
var wunderCtlOnce sync.Once

func GetWunderCtlClient() wunderClient {
	rootContext := context.Background()
	wunderCtlOnce.Do(func() {
		onceCtx, onceCancel := context.WithCancel(rootContext)
		wunderCtlClient = wunderClient{
			WunderCtlClient: commands.NewClient(),
			RootContext:     rootContext,
			Context: Context{
				Ctx:    onceCtx,
				Cancel: onceCancel,
			},
		}
	})
	if wunderCtlClient.Context.Ctx.Err() != nil {
		if wunderCtlClient.Context.Ctx.Err().Error() == "context canceled" {
			onceCtx, onceCancel := context.WithCancel(rootContext)
			wunderCtlClient.Context = Context{
				Ctx:    onceCtx,
				Cancel: onceCancel,
			}
		}
	}

	return wunderCtlClient
}

func Stop() {
	wdgClient := GetWunderCtlClient()
	wdgClient.Context.Cancel()
}

func GetWdgConfig() domain.SystemConfig {
	setting := domain.Setting{}
	content, err := ioutil.ReadFile(utils.GetSettingPath())
	if err != nil {
		panic(fmt.Sprintf("init wundergraph err=%+v", err))
	}
	if len(content) != 0 {
		err = json.Unmarshal(content, &setting)
		if err != nil {
			panic(fmt.Sprintf("init wundergraph err=%+v", err))
		}
	}
	if setting.System.DevSwitch == true {
		setting.System.ForcedJumpSwitch = true
	}
	return setting.System
}

func GetWdgArgs(config domain.SystemConfig) []byte {
	args := struct {
		LogLevel                   string `json:"logLevel"`
		ExcludeServer              bool   `json:"excludeServer"`
		ListenAddr                 string `json:"listenAddr"`                 // 默认9991端口 ,可命令行更改
		MiddlewareListenPort       string `json:"middlewareListenPort"`       // 默认9992端口 ,可命令行更改
		EnableDebugMode            bool   `json:"enableDebugMode"`            // debug开关,命令行获取，默认关
		DisableForceHttpsRedirects bool   `json:"disableForceHttpsRedirects"` //  https强制跳转,命令行获取,默认开
		EnableIntrospection        bool   `json:"enableIntrospection"`        //  内省开关,命令行获取,默认不开启,
		LogPath                    string `json:"logPath"`                    //  日志文件写入路径
	}{
		LogLevel:                   config.LogLevel,
		ListenAddr:                 fmt.Sprintf("%s:%s", utils.GetFireBoomHost(), config.ApiPort),
		MiddlewareListenPort:       config.MiddlewarePort,
		EnableDebugMode:            config.DevSwitch,
		DisableForceHttpsRedirects: config.ForcedJumpSwitch,
		LogPath:                    utils.GetWdgLogPath(),
	}

	// 初始化配置
	argContent, _ := json.Marshal(args)
	return argContent
}

func WdgGenerate(wdgClient wunderClient, args []byte, WdgLogChan, HookLogChan chan string) error {
	return wdgClient.WunderCtlClient.Generate(wdgClient.Context.Ctx, args, WdgLogChan, HookLogChan)
}

func WdgStart(wdgClient wunderClient, args []byte, WdgLogChan, HookLogChan chan string) error {
	return wdgClient.WunderCtlClient.WdgStart(wdgClient.Context.Ctx, args, WdgLogChan, HookLogChan)
}

func KillPid() {
	// 等于-1说明端口未找到，或者未启动
	pid := PortInUse(9992)
	if pid == -1 {
		return
	}
	var outBytes bytes.Buffer
	cmdStr := fmt.Sprintf("kill -9  %d", pid)
	cmd := exec.Command("cmd", "-k", cmdStr)
	cmd.Stdout = &outBytes
	cmd.Run()

	return
}

func PortInUse(portNumber int) int {
	res := -1
	var outBytes bytes.Buffer
	cmdStr := fmt.Sprintf("netstat -ano -p tcp | findstr %d", portNumber)
	cmd := exec.Command("bash", "-p", cmdStr)
	cmd.Stdout = &outBytes
	cmd.Run()
	resStr := outBytes.String()
	r := regexp.MustCompile(`\s\d+\s`).FindAllString(resStr, -1)
	if len(r) > 0 {
		pid, err := strconv.Atoi(strings.TrimSpace(r[0]))
		if err != nil {
			res = -1
		} else {
			res = pid
		}
	}
	return res
}
