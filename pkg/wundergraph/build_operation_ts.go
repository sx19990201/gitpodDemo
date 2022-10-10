package wundergraph

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fire_boom/domain"
	"github.com/fire_boom/utils"
	"github.com/labstack/gommon/log"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

type operationData struct {
	OperationConfig string
}

func (w *wdg) buildOperationData(ctx context.Context) (operationData, error) {
	od := operationData{}
	// 全局operations设置
	w.buildGlobalOperationContent(ctx, &od)
	return od, nil
}

func (w *wdg) buildGlobalOperationContent(ctx context.Context, cd *operationData) {
	configStr := GetWdgOperGlobalConfig()
	// 获取operation详情
	operations, err := w.repository.or.FindOperations(ctx)
	if err != nil {
		log.Error("buildGlobalOperationContent FindOperations err : ", err)
	}
	customConfig := fmt.Sprintf(`custom: {`)
	for _, row := range operations {
		// 如果没有开启则不生成
		if row.Enable == domain.SwitchOff {
			continue
		}
		operName := strings.Split(filepath.Base(row.Path), ".")[0]
		// 解析这个配置
		path := fmt.Sprintf("%s%s", utils.GetOperationSettingPath(), row.Path)
		// 判断文件是否存在
		if !utils.FileExist(path) {
			continue
		}
		content, err := ioutil.ReadFile(path)
		if err != nil {
			log.Error(fmt.Sprintf("buildGlobalOperationContent get %s operations content fail , err : ", row.Path), err)
		}
		var operSetting domain.OperationSetting
		// 如果内容为空，说明没有配置
		if len(content) == 0 {
			continue
		}
		err = json.Unmarshal(content, &operSetting)
		if err != nil {
			log.Error(fmt.Sprintf("buildGlobalOperationContent get %s operations json.Unmarshal fail , err : ", row.Path), err)
		}
		var operAuth WdgAuthentication
		var operCaching WdgCaching
		var operLiveQuery WdgLiveQuery
		operAuth.Required = operSetting.AuthenticationRequired
		// 查询 需要添加缓存等配置，修改则不用，只有身份验证的配置
		authStr := ""

		authStr = operAuth.GetWdgConfig()

		if OperationTypeMap[row.OperationType] == queries {
			// 如果是查询，判断设置是否启用
			operCaching.Enable = operSetting.CachingEnable
			operCaching.StaleWhileRevalidate = operSetting.CachingStaleWhileRevalidate
			operCaching.MaxAge = operSetting.CachingMaxAge

			operLiveQuery.PollingIntervalSeconds = operSetting.LiveQueryPollingIntervalSeconds
			operLiveQuery.Enable = operSetting.LiveQueryEnable

			cachingStr := operCaching.GetWdgConfig()
			liveQueryStr := operLiveQuery.GetWdgConfig()

			customConfig = fmt.Sprintf(`%s%s : config => ({
				...config,
				%s 
				%s 
				%s 
			}), `, customConfig, operName, authStr, cachingStr, liveQueryStr)
			continue
		}
		customConfig = fmt.Sprintf(`%s%s : config => ({
				...config,
				%s 
			}), `, customConfig, operName, authStr)
	}
	customConfig = fmt.Sprintf(`%s },`, customConfig)
	configStr = fmt.Sprintf(`%s %s },`, configStr, customConfig)
	reg := regexp.MustCompile("\"([a-zA-Z]\\w*)\":")
	configStr = reg.ReplaceAllString(configStr, `$1:`)
	cd.OperationConfig = configStr
}

func GetWdgOperGlobalConfig() string {
	// 获取全局配置
	settingByte, err := ioutil.ReadFile(utils.GetOperationGlobalSettingPath())
	if err != nil {
		log.Error("Failed to get global settings , err : ", err)
	}
	var globalSetting domain.OperationSetting
	err = json.Unmarshal(settingByte, &globalSetting)
	if err != nil {
		log.Error("buildGlobalOperationContent json.Unmarshal  globalSetting err : ", err)
	}

	globalAuth := WdgAuthentication{
		Required: globalSetting.AuthenticationRequired,
	}
	globalDefaultConfig := WdgDefaultConfig{
		AuthConfig: globalAuth,
	}
	globalOperationsConfig := WdgOperationsConfig{
		DefaultConfig: globalDefaultConfig,
	}
	// 设置默认值
	if globalSetting.CachingStaleWhileRevalidate == 0 {
		globalSetting.CachingStaleWhileRevalidate = 60
	}
	if globalSetting.CachingMaxAge == 0 {
		globalSetting.CachingMaxAge = 60
	}
	globalCache := WdgCaching{
		Enable:               globalSetting.CachingEnable,
		StaleWhileRevalidate: globalSetting.CachingStaleWhileRevalidate,
		MaxAge:               globalSetting.CachingMaxAge,
	}
	globalLiveQuery := WdgLiveQuery{
		Enable:                 globalSetting.LiveQueryEnable,
		PollingIntervalSeconds: globalSetting.LiveQueryPollingIntervalSeconds,
	}
	globalCacheBytes, _ := json.Marshal(globalCache)
	cacheStr := fmt.Sprintf("caching : %s,", string(globalCacheBytes))

	globalLiveQueryBytes, _ := json.Marshal(globalLiveQuery)
	liveQueryStr := fmt.Sprintf("liveQuery: %s,", string(globalLiveQueryBytes))
	globalOperationsConfigBytes, _ := json.Marshal(globalOperationsConfig)
	defaultStr := string(globalOperationsConfigBytes)
	defaultStr = strings.TrimLeft(defaultStr, "{}")
	defaultStr = strings.ReplaceAll(defaultStr, "}}}", "}}")

	defaultStr = fmt.Sprintf("%s,", defaultStr)
	// 处理双引号
	reg := regexp.MustCompile("\"([a-zA-Z]\\w*)\":")
	cacheStr = reg.ReplaceAllString(cacheStr, `$1:`)
	liveQueryStr = reg.ReplaceAllString(liveQueryStr, `$1:`)
	defaultStr = reg.ReplaceAllString(defaultStr, `$1:`)

	queriesContent := fmt.Sprintf(`%s : (config) => ({
			...config,
			%s
			%s
	}),`, queries, cacheStr, liveQueryStr)
	// 修改的全局配置暂时没有
	mutationsContent := fmt.Sprintf(`%s : (config) => ({
		...config,
	}),`, mutations)
	// 订阅的全局配置暂时没有
	subscriptionsContent := fmt.Sprintf(`%s : (config) => ({
		...config,
	}),`, subscriptions)

	operationsConfig := fmt.Sprintf(`operations: { 
		%s
		%s
 		%s
       	%s
	`, defaultStr, queriesContent, mutationsContent, subscriptionsContent)

	return operationsConfig
}
