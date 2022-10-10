package http

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/fire_boom/pkg/wdgfunc"
	"github.com/fire_boom/pkg/wundergraph"
	"github.com/fire_boom/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (
	Building = "编译中"
	Started  = "已启动"
	Stop     = "已关闭"
)

type StatusBarResponse struct {
	EngineStatus string    `json:"engineStatus"` // 引擎状态
	HookStatus   string    `json:"hookStatus"`   // 钩子状态
	ErrorInfo    ErrorInfo `json:"errorInfo"`    // 错误信息
	Env          string    `json:"env"`          // 环境
}

type BarStaticResponse struct {
	Version string `json:"version"`
	Env     string `json:"env"`
}

type ErrorInfo struct {
	WarnTotal int64 `json:"warnTotal"`
	ErrTotal  int64 `json:"errTotal"`
}

// WdgStatus 默认为已关闭
var WdgStatus = Stop

type WdgHandler struct {
	db *sql.DB
}

func InitWdgHandler(e *echo.Echo, db *sql.DB) {
	NewWdgHandler(e, db)
}

func NewWdgHandler(e *echo.Echo, db *sql.DB) {
	handler := &WdgHandler{db: db}
	v1 := e.Group("/api/v1")
	{
		wdg := v1.Group("/wdg")
		{
			wdg.GET("/close", handler.Stop)
			wdg.GET("/state", echo.WrapHandler(http.HandlerFunc(State)))
			wdg.GET("/log", echo.WrapHandler(http.HandlerFunc(GetLog)))
			wdg.GET("/start", handler.Start)
			wdg.GET("/reStart", handler.ReStart)
			wdg.GET("/barOnce", handler.BarOnce)
		}
	}
}

type WriteDetail struct {
	ID    uuid.UUID
	Write http.ResponseWriter
}

var StartCount = 0

// BarOnce @Title BarOnce
// @Description
// @Accept  json
// @Tags  wdg
// @Success 200 "重启成功"
// @Failure 400	"重启失败"
// @Router /api/v1/wdg/barOnce  [GET]
func (w *WdgHandler) BarOnce(c echo.Context) error {
	StartCount++
	// 清空日志
	setting, err := utils.GetSettingConfig()
	if err != nil {
		log.Error("get setting config fail err : ", err)
	}
	env := "生产模式"
	if setting.System.DevSwitch {
		env = "开发模式"
	}
	version := setting.Version.VersionNum
	return c.JSON(http.StatusOK, SuccessWriteResult(BarStaticResponse{
		Env:     env,
		Version: version,
	}))
}

// GetLog @Title GetLog
// @Description 查看日志
// @Accept  json
// @Tags  wdg
// @Param type formData integer true "类型 1-wdg日志（9991），2-hook日志（9992）"
// @Success 200 "查看日志成功"
// @Failure 400	"查看日志失败"
// @Router /api/v1/wdg/log  [GET]
func GetLog(w http.ResponseWriter, r *http.Request) {
	// 这一步很重要,要text类型的context-type才支持flush
	w.Header().Set("Content-Type", "text/html;charset=utf8")

	resultRows := make([]LogRow, 0)
	// 获取wdg日志
	resultRows = append(resultRows, GetTailLogs(utils.ReadFileLine(utils.GetWdgLogPath()), WdgLogType)...)
	// 获取hooks日志
	resultRows = append(resultRows, GetTailLogs(utils.ReadFileLine(utils.GetHookLogPath()), HookLogType)...)

	for _, row := range resultRows {
		f, _ := w.(http.Flusher)
		content, _ := json.Marshal(row)
		fmt.Fprintf(w, "%s \n", string(content))
		time.Sleep(time.Millisecond * 20)
		f.Flush()
	}
	id := uuid.New().String()
	FlusherMsgCMap.Set(id, FlusherMsg{
		FlusherType: LogTypeFlag,
		MsgWriter:   w,
	})
	<-r.Context().Done()
	FlusherMsgCMap.Remove(id)
}

func GetTailLogs(arr []string, logType int64) (result []LogRow) {
	if len(arr) > 10 {
		arr = arr[len(arr)-10:]
	}

	for i := 0; i < len(arr); i++ {
		row := arr[i]
		logMap := make(map[string]interface{}, 0)
		json.Unmarshal([]byte(row), &logMap)
		if len(logMap) == 0 {
			continue
		}
		result = append(result, LogRow{
			Time:    utils.InterfaceToString(logMap["time"]),
			Level:   utils.InterfaceToString(logMap["level"]),
			Msg:     utils.InterfaceToString(logMap["msg"]),
			LogType: logType,
		})
	}

	return
}

// State @Title State
// @Description 查看状态
// @Accept  json
// @Tags  wdg
// @Success 200 "查看状态成功"
// @Failure 400	"查看状态失败"
// @Router /api/v1/wdg/state  [GET]
func State(w http.ResponseWriter, r *http.Request) {
	// 这一步很重要,要text类型的context-type才支持flush
	w.Header().Set("Content-Type", "text/html;charset=utf8")
	f, _ := w.(http.Flusher)
	statusBar := StatusBarResponse{
		EngineStatus: WdgStatus,
		HookStatus:   WdgStatus,
		ErrorInfo: ErrorInfo{
			WarnTotal: 1,
			ErrTotal:  1,
		},
	}
	result, _ := json.Marshal(statusBar)
	fmt.Fprintf(w, "%s \n", result)
	// 刷新,推送给前端
	f.Flush()

	id := uuid.New().String()
	FlusherMsgCMap.Set(id, FlusherMsg{
		FlusherType: 0,
		MsgWriter:   w,
	})

	<-r.Context().Done()
	// 请求结束删掉该
	//WriteSyncMap.Del(id)
	//FlusherMsgMap.Del(id)
	FlusherMsgCMap.Remove(id)
}

// Stop @Title stop
// @Description 停止
// @Accept  json
// @Tags  wdg
// @Success 200 "停止成功"
// @Failure 400	"停止失败"
// @Router /api/v1/wdg/close  [GET]
func (w *WdgHandler) Stop(c echo.Context) error {
	wdgClient := wdgfunc.GetWunderCtlClient()
	wdgClient.Context.Cancel()

	WdgStatus = Stop
	StatusChan <- WdgStatus
	return c.JSON(http.StatusOK, SuccessResult())
}

// ReStart @Title ReStart
// @Description 重启
// @Accept  json
// @Tags  wdg
// @Success 200 "重启成功"
// @Failure 400	"重启失败"
// @Router /api/v1/wdg/reStart  [GET]
func (w *WdgHandler) ReStart(c echo.Context) error {
	StartCount++
	// 清空日志
	utils.WriteFile(utils.GetWdgLogPath(), "")
	utils.WriteFile(utils.GetHookLogPath(), "")
	go w.RestartWdg()

	return c.JSON(http.StatusOK, SuccessResult())
}

// Start @Title Start
// @Description 启动
// @Accept  json
// @Tags  wdg
// @Success 200 "启动成功"
// @Failure 400	"启动失败"
// @Router /api/v1/wdg/start  [GET]
func (w *WdgHandler) Start(c echo.Context) error {
	StartCount++
	w.ReloadConfig()
	// 清空日志
	utils.WriteFile(utils.GetWdgLogPath(), "")
	utils.WriteFile(utils.GetHookLogPath(), "")
	go StartWdg()
	return c.JSON(http.StatusOK, SuccessResult())
}

func (w *WdgHandler) RestartWdg() {
	// 先停止启动
	wdgfunc.Stop()
	WdgStatus = Stop
	StatusChan <- WdgStatus

	// 再启动
	w.ReloadConfig()
	go StartWdg()
}

func (w *WdgHandler) ReloadConfig() {

	//合成文件
	err := wundergraph.ReloadWdgFile(w.db)
	// 合成失败，重试一次
	if err != nil {
		wundergraph.ReloadWdgFile(w.db)
	}
	WdgStatus = Building
	StatusChan <- WdgStatus
}

func StartWdg() {
	config := wdgfunc.GetWdgConfig()
	args := wdgfunc.GetWdgArgs(config)

	wdgClient := wdgfunc.GetWunderCtlClient()
	// 获取wdgclient
	err := wdgfunc.WdgGenerate(wdgClient, args, WdgLogChan, HookLogChan)
	if err != nil {
		// 发送启动失败消息
		WdgStatus = Stop
		StatusChan <- WdgStatus
		log.Error("wdg generate fail ,err : ", err)
	}

	wdgfunc.KillPid()
	WdgStatus = Started
	//PushMsg()
	StatusChan <- WdgStatus

	if err := wdgfunc.WdgStart(wdgClient, args, WdgLogChan, HookLogChan); err != nil {
		// 发送启动失败消息
		WdgStatus = Stop
		StatusChan <- WdgStatus
		log.Error("wdg start fail ,err : ", err)
	}
}
