package wundergraph

import "fmt"

type WdgOperationsConfig struct {
	DefaultConfig WdgDefaultConfig `json:"defaultConfig"`
}

// WdgDefaultConfig 全局配置
type WdgDefaultConfig struct {
	AuthConfig WdgAuthentication `json:"authentication"`
}

type WdgAuthentication struct {
	Required bool `json:"required"`
}

func (w *WdgAuthentication) GetWdgConfig() (result string) {
	if w == nil {
		return ""
	}

	if w.Required == false {
		return result

	}
	required := fmt.Sprintf("required: %v ,", w.Required)
	result = fmt.Sprintf(`authentication:{
		%s
	},`, required)
	return result
}

type WdgGlobalConfig struct {
	Caching   WdgCaching   `json:"caching"`
	LiveQuery WdgLiveQuery `json:"liveQuery"`
}

type WdgCaching struct {
	Enable               bool  `json:"enable"`
	StaleWhileRevalidate int64 `json:"staleWhileRevalidate,default=60"`
	MaxAge               int64 `json:"maxAge,default=60"`
	Public               bool  `json:"public,default=true"`
}

func (w *WdgCaching) GetWdgConfig() string {
	if w == nil {
		return ""
	}

	staleWhileRevalidate := ""
	maxAge := ""
	enable := ""
	if w.StaleWhileRevalidate != 0 {
		staleWhileRevalidate = fmt.Sprintf("staleWhileRevalidate: %v ,", w.StaleWhileRevalidate)
	}
	if w.MaxAge != 0 {
		maxAge = fmt.Sprintf("maxAge: %v ,", w.MaxAge)
	}
	if w.Enable != false {
		enable = fmt.Sprintf("enable: true ,")
	}

	result := fmt.Sprintf(`caching:{
		...config.caching,
		%s
		%s
		%s
	},`, staleWhileRevalidate, maxAge, enable)
	return result
}

type WdgLiveQuery struct {
	Enable                 bool  `json:"enable,default=true"`
	PollingIntervalSeconds int64 `json:"pollingIntervalSeconds,default=1"`
}

func (w *WdgLiveQuery) GetWdgConfig() string {
	if w == nil {
		return ""
	}

	pollingIntervalSeconds := ""
	enable := ""
	if w.PollingIntervalSeconds != 0 {
		pollingIntervalSeconds = fmt.Sprintf("pollingIntervalSeconds: %v ,", w.PollingIntervalSeconds)
	}
	if w.Enable != false {
		enable = fmt.Sprintf("enable: true ,")
	}

	result := fmt.Sprintf(`liveQuery:{
		...config.liveQuery,
		%s
		%s
	},`, pollingIntervalSeconds, enable)
	return result
}
