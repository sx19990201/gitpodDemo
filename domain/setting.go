package domain

const (
	DevSwitchOn           = "0"
	DevSwitchOff          = "1"
	ForcedJumpSwitchOn    = "0"
	ForcedJumpSwitchOff   = "1"
	LogLevelInfo          = "0"
	LogLevelDebug         = "1"
	LogLevelError         = "2"
	infoName              = "info"
	debugName             = "debug"
	errorName             = "error"
	DefaultApiPort        = "9991"
	DefaultMiddlewarePort = "9992"
)

type Setting struct {
	System      SystemConfig      `json:"system"`      // 系统
	Version     VersionConfig     `json:"version"`     // 版本
	Environment EnvironmentConfig `json:"environment"` // 环境变量
}

type SystemConfig struct {
	ApiPort          string `json:"apiPort"`          // api端口
	MiddlewarePort   string `json:"middlewarePort"`   // 中间件端口
	LogLevel         string `json:"logLevel"`         // 日志水平 0-Info 1-debug 2-error
	DevSwitch        bool   `json:"devSwitch"`        // 开发者模式开关 true开 false关 0-开(debug模式 用up 启动) 1-关(生产模式，用start启动)
	ForcedJumpSwitch bool   `json:"forcedJumpSwitch"` // 强制跳转 强制重定向跳转，开启后强制使用https协议
	DebugSwitch      bool   `json:"debugSwitch"`      // 调试开关 0-关 1-开
	//EnvType          string `json:"envType"`          // 开发环境
}

func (s *SystemConfig) GetDebugCMD() bool {
	return s.DebugSwitch
}

func (s *SystemConfig) GetApiPortCMD() string {
	if s.ApiPort == "" {
		s.ApiPort = DefaultApiPort
	}
	return s.ApiPort
}

func (s *SystemConfig) GetMiddlewarePortCMD() string {
	if s.MiddlewarePort == "" {
		s.MiddlewarePort = DefaultMiddlewarePort
	}
	return s.MiddlewarePort
}

func (s *SystemConfig) GetLogLevelCMD() string {
	if s.LogLevel == LogLevelInfo {
		return infoName
	}
	if s.LogLevel == LogLevelDebug {
		return debugName
	}
	if s.LogLevel == LogLevelError {
		return errorName
	}
	return ""
}

func (s *SystemConfig) GetDevSwitchCMD() string {
	if s.DevSwitch == true {
		return "upsdk"
	}
	if s.DevSwitch == false {
		return "start"
	}
	return "up"
}

func (s *SystemConfig) GetForcedJumpCMD() string {
	if s.ForcedJumpSwitch == false {
		return "--disable-force-https-redirects "
	}
	return ""
}

type VersionConfig struct {
	VersionNum    string `json:"versionNum"`    // 版本号
	PrismaVersion string `json:"prismaVersion"` // prisma版本
	Copyright     string `json:"copyright"`     // 版权
}

type EnvironmentConfig struct {
	EnvironmentList []EnvironmentDetail `json:"environmentList"` // 环境变量列表
	SystemVariable  string              `json:"systemVariable"`  // 系统变量
}

type EnvironmentDetail struct {
	Name string `json:"name"` // 变量名
	Dev  string `json:"dev"`  // 开发环境
	Pro  string `json:"pro"`  // 生产环境
}
