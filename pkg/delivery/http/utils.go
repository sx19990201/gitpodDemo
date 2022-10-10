package http

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/fire_boom/utils"
	cmap "github.com/orcaman/concurrent-map/v2"
	"net/http"
	"os"
	"sync"
)

const (
	StateTypeFlag = iota
	LogTypeFlag
)
const (
	WdgLogType  = 1
	HookLogType = 2
)

var WdgLogChan = make(chan string, 1024)
var HookLogChan = make(chan string, 1024)
var StatusChan = make(chan string, 1024)

//var LogChan = make(chan LogRow, 1024)

type LogRow struct {
	Time    string `json:"time"`
	Level   string `json:"level"`
	Msg     string `json:"msg"`
	LogType int64  `json:"logType"`
}

var Msg http.ResponseWriter

var FlusherMsgCMap = cmap.New[FlusherMsg]()
var FlusherMsgMap = flusherMsgMap{
	FlusherMap: make(map[string]FlusherMsg, 0),
	mutex:      &sync.RWMutex{},
}

type FlusherMsg struct {
	FlusherType int64               // 类型 0-状态 1-日志
	MsgWriter   http.ResponseWriter // 数据
}

type flusherMsgMap struct {
	FlusherMap map[string]FlusherMsg
	mutex      *sync.RWMutex
}

func (f *flusherMsgMap) Get(id string) FlusherMsg {
	f.mutex.RLock()

	defer f.mutex.RUnlock()
	return f.FlusherMap[id]
}

func (f *flusherMsgMap) Store(id string, val FlusherMsg) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.FlusherMap[id] = val
	return
}

func (f *flusherMsgMap) Del(id string) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	delete(f.FlusherMap, id)
	return
}

// Keys 所有的key
func (f flusherMsgMap) Keys() []string {
	f.mutex.RLock() // 加锁(读锁定)
	keys := make([]string, len(f.FlusherMap))
	for k := range f.FlusherMap {
		keys = append(keys, k)
	}
	f.mutex.RUnlock() // 解锁
	return keys
}

func PushLogMsg() {
	go func() {
		for msg := range WdgLogChan {
			FlushLogMsg(WdgLogType, msg)
		}
	}()
	go func() {
		for msg := range HookLogChan {
			FlushLogMsg(HookLogType, msg)
		}
	}()
}

func PushStatusMsg() {
	go func() {
		for msg := range StatusChan {
			FlushStatusMsg(msg)
		}
	}()
}

func FlushStatusMsg(content string) {
	statusBar := StatusBarResponse{
		EngineStatus: content,
		HookStatus:   content,
		ErrorInfo: ErrorInfo{
			WarnTotal: 1,
			ErrTotal:  1,
		},
	}
	msg, _ := json.Marshal(statusBar)
	FlusherMap(StateTypeFlag, string(msg))
}

func FlushLogMsg(logType int64, content string) {
	if logType == WdgLogType {
		FileAppendWrite(utils.GetWdgLogPath(), content)
		FlusherMap(LogTypeFlag, content)
	}
	if logType == HookLogType {
		FileAppendWrite(utils.GetHookLogPath(), content)
		FlusherMap(LogTypeFlag, content)
	}
}

func FlusherMap(flusherType int64, content string) {
	for item := range FlusherMsgCMap.Iter() {
		flusher := item.Val

		if flusher.FlusherType != flusherType {
			continue
		}

		Flush(flusher.MsgWriter, content)
	}
}

// FileAppendWrite 追加写入
func FileAppendWrite(filePath, msg string) {
	// 文件不存在
	if !utils.FileExist(filePath) {
		// 创建
		utils.WriteFile(filePath, "")
	}

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	defer file.Close()
	//写入文件时，使用带缓存的 *Writer
	logWrite := bufio.NewWriter(file)
	logWrite.WriteString(fmt.Sprintf("%s \n", msg))
	//Flush将缓存的文件真正写入到文件中
	logWrite.Flush()
}

func Flush(write http.ResponseWriter, msg string) {
	write.Header().Set("Content-Type", "text/html;charset=utf8")
	// 断言flusher
	f, _ := write.(http.Flusher)
	// 内容写入responseWrite
	fmt.Fprintf(write, "%s", msg)
	// 刷新,推送给前端
	f.Flush()
}
