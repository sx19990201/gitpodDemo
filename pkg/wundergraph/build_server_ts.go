package wundergraph

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fire_boom/domain"
	"github.com/fire_boom/pkg/usecase"
	"github.com/fire_boom/utils"
	"github.com/labstack/gommon/log"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	wdgAuthHookPath       = "./hooks/auth"
	wdgGlobalHookPath     = "./hooks/global"
	wdgMockHookPath       = "./hooks/mock"
	wdgCustomizeHookPath  = "./hooks/customize"
	wdgOperationsHookPath = "./hooks/operations"
	mockResolveName       = "mockResolveName"
)
const defaultServer = `{
			apiNamespace: "gql",
			serverName: "gql",
			schema: new GraphQLSchema({
				query: new GraphQLObjectType({
					name: 'RootQueryType',
					fields: {
						hello: {
							type: GraphQLString,
							resolve() {
								return 'world';
							},
						},
					},
				}),
			}),
		}`

type serverData struct {
	ImportPath  string
	HookContent string
}

func (w *wdg) buildServerData(ctx context.Context) (serverData, error) {
	sd := serverData{}
	// 全局operations设置
	w.buildServerContent(ctx, &sd)
	return sd, nil
}

func (w *wdg) buildServerContent(ctx context.Context, cd *serverData) {
	timeoutContext := time.Duration(utils.GetTimeOut()) * time.Second
	hookUC := usecase.NewHooksUseCase(timeoutContext)
	// 身份鉴权hooks postAuthentication 和 mutatingPostAuthentication
	authHookMap := GetAuthHooks(hookUC)
	fmt.Println(authHookMap)
	// 获取全局钩子
	globalHookMap := GetGlobalHooks(hookUC)
	fmt.Println(globalHookMap)
	operations, err := w.repository.or.FindOperations(ctx)
	if err != nil {
		// TODO
		log.Error("")
	}
	operationMap := make(map[string]domain.FbOperations)
	for _, row := range operations {
		key := strings.Split(row.Path, ".")[0]
		// 类型映射
		row.OperationType = OperationTypeMap[row.OperationType]
		operationMap[key] = row
	}
	operationHooks := GetOperationHoos(hookUC, operationMap)
	fmt.Println(operationHooks)
	customHooks := GetGraphqlServers(ctx, w.repository.dsr)
	fmt.Println(customHooks)

	globalImportName, globalConfig := GetGlobalConfig(globalHookMap)
	authImportName, authConfig := GetAuthConfig(authHookMap)
	operationsImportName, operationsConfig := GetOperationConfig(operationHooks)
	customImportName, customConfig := GetGraphqlServersConfig(customHooks)

	importConfig := fmt.Sprintf("%s %s %s %s",
		globalImportName, authImportName, operationsImportName, customImportName)
	hooksConfig := fmt.Sprintf("hooks:{ %s %s %s },", globalConfig, authConfig, operationsConfig)
	contentConfig := fmt.Sprintf("%s %s ", hooksConfig, customConfig)

	cd.ImportPath = importConfig
	cd.HookContent = contentConfig
}

func GetGlobalHooks(hookUC *usecase.HooksUseCase) map[string]string {
	// 全局hooks onRequest onResponse
	result := make(map[string]string, 0)
	globalHooks, err := hookUC.FindHooksByPath(utils.GetGlobalHookPathPrefix(), "", domain.GlobalHook)
	if err != nil {
		log.Error("")
	}
	for _, hook := range globalHooks {
		// 如果hooks关闭则跳过
		if hook.HookSwitch == utils.OFF {
			continue
		}
		path := fmt.Sprintf("%s/%s", wdgGlobalHookPath, hook.HookName)
		importPath := fmt.Sprintf("import %s from '%s'; \n", hook.HookName, path)
		result[hook.HookName] = importPath
	}
	return result
}

func GetGlobalConfig(hooksMap map[string]string) (string, string) {
	importConfig := ""
	onRequestConfig := ""
	onResponseConfig := ""
	for importName, importPath := range hooksMap {
		importConfig = fmt.Sprintf("%s %s", importConfig, importPath)
		if importName == "onRequest" {
			onRequestConfig = fmt.Sprintf(`onRequest:{
				hook: %s,
			},`, importName)
			continue
		}
		if importName == "onResponse" {
			onResponseConfig = fmt.Sprintf(`onResponse:{
				hook: %s,
			},`, importName)
			continue
		}
	}
	globalConfig := fmt.Sprintf(`global: {
                httpTransport:{
                    %s
                    %s
                }
            },`, onRequestConfig, onResponseConfig)
	if onRequestConfig == "" && onResponseConfig == "" {
		return "", ""
	}
	return importConfig, globalConfig
}

func GetAuthHooks(hookUC *usecase.HooksUseCase) map[string]string {
	// ../../static/hooks/auth/mutatingPostAuthentication
	result := make(map[string]string, 0)
	authHooks, err := hookUC.FindHooksByPath(utils.GetAuthGlobalHookPathPrefix(), "", domain.AuthGlobalHook)
	if err != nil {
		log.Error("get hook path error, err : ", err)
	}
	for _, hook := range authHooks {
		// 如果hooks关闭则跳过
		if hook.HookSwitch == utils.OFF {
			continue
		}
		path := fmt.Sprintf("%s/%s", wdgAuthHookPath, hook.HookName)
		importPath := fmt.Sprintf("import %s from '%s'; \n", hook.HookName, path)
		result[hook.HookName] = importPath
	}
	return result
}

func GetAuthConfig(hooksMap map[string]string) (string, string) {
	importConfig := ""
	authConfig := ""
	contentConfig := ""
	for importName, importPath := range hooksMap {
		importConfig = fmt.Sprintf("%s %s", importConfig, importPath)
		contentConfig = fmt.Sprintf("%s %s,", contentConfig, importName)
	}
	authConfig = fmt.Sprintf(`authentication: {
				%s
            },`, contentConfig)
	return importConfig, authConfig
}

func GetOperationHoos(hookUC *usecase.HooksUseCase, operationMap map[string]domain.FbOperations) map[string]map[string]map[string]string {
	hookFilePaths := utils.GetDicAllChildFilePath(utils.GetOperationsHookPathPrefix())
	result := make(map[string]map[string]map[string]string, 0)

	for _, path := range hookFilePaths {
		operMap := make(map[string]map[string]string, 0)
		hookMap := make(map[string]string)

		// 处理路径
		rePath := filepath.ToSlash(utils.GetOperationsHookPathPrefix())
		path = filepath.ToSlash(path)
		path = strings.ReplaceAll(path, rePath, "")
		// 找到该operation
		operation, ok := operationMap[strings.Split(path, "_")[0]]
		// 数据库没有该文件
		if !ok {
			// TODO 可以收集这些文件做个通知，让数据同步
			continue
		}
		// 如果该operations关闭，则跳过
		if operation.Enable == domain.SwitchOff {
			continue
		}
		// 获取该operation的钩子列表
		operationHooks, err := hookUC.FindHooksByPath(utils.GetOperationsHookPathPrefix(), strings.Split(path, "_")[0], domain.OperationHook)
		if err != nil {
			log.Error(err)
		}
		// 遍历钩子列表
		for _, hook := range operationHooks {
			// 如果按钮被关闭则跳过
			if hook.HookSwitch == utils.OFF {
				continue
			}
			hookName := filepath.Base(fmt.Sprintf("%s_%s", operation.Path, hook.HookName))
			hookPath := fmt.Sprintf("%s/%s", wdgOperationsHookPath, hookName)
			importPath := fmt.Sprintf("import %s from '%s'; \n", hookName, hookPath)
			hookMap[hookName] = importPath
		}
		//  只读取开启的mock
		mockPath := fmt.Sprintf("%s%s%s%s", utils.GetMockPath(), strings.Split(path, "_")[0], utils.GetHooksSuffix(), utils.GetSwitchState(utils.ON))
		// 如果mock存在
		if utils.FileExist(mockPath) {
			// 如果按钮开启了则加入map
			if !strings.Contains(mockPath, utils.GetSwitchState(utils.OFF)) {
				hookName := filepath.Base(fmt.Sprintf("%s_%s", strings.Split(path, "_")[0], mockResolveName))
				importMockPath := fmt.Sprintf("%s/%s", wdgMockHookPath, hookName)

				importPath := fmt.Sprintf("import %s from '%s'; \n", hookName, strings.Split(importMockPath, "_")[0])
				hookMap[hookName] = importPath
			}
		}
		key := filepath.Base(strings.Split(path, "_")[0])
		operMap[key] = hookMap
		result[operation.OperationType] = operMap
	}
	return result
}

func GetOperationConfig(hookMap map[string]map[string]map[string]string) (string, string) {
	config := ""
	importConfig := ""
	methodConfig := ""
	for method, operationMap := range hookMap {
		operationConfig := ""
		for operationName, operationHooksMap := range operationMap {
			operationConfig = fmt.Sprintf(`%s%s :{ `, operationConfig, operationName)
			for hookName, importPath := range operationHooksMap {
				importConfig = fmt.Sprintf("%s %s", importConfig, importPath)
				hookNameConfig := strings.Split(hookName, "_")[1]
				operationConfig = fmt.Sprintf("%s %s:%s,", operationConfig, hookNameConfig, hookName)
			}
			operationConfig = fmt.Sprintf(`%s },`, operationConfig)
		}
		methodConfig = fmt.Sprintf(`%s:{%s},`, method, operationConfig)
		config = fmt.Sprintf("%s%s", config, methodConfig)
	}
	return importConfig, config
}

func GetGraphqlServers(ctx context.Context, dsr domain.DataSourceRepository) map[string]string {
	dataSources, err := dsr.FindDataSources(ctx)
	if err != nil {
		// TODO
	}
	result := make(map[string]string)
	for _, row := range dataSources {
		if row.SourceType != domain.SourceTypeCustomize {
			continue
		}
		// 开关关闭则跳过
		if row.Switch == domain.SwitchOff {
			continue
		}
		config := domain.CustomizeConfig{}
		err := json.Unmarshal([]byte(row.Config), &config)
		if err != nil {
			// TODO 日志
		}
		config.Schema = config.ApiNamespace
		path := fmt.Sprintf("%s/%s", wdgCustomizeHookPath, config.ApiNamespace)
		importPath := fmt.Sprintf("import %s from '%s'; \n", config.ApiNamespace, path)
		val, err := json.Marshal(config)
		if err != nil {
			// TODO
		}
		result[importPath] = string(val)
	}
	return result
}

func GetGraphqlServersConfig(cusMap map[string]string) (string, string) {
	config := ""
	importConfig := ""
	cusConfig := ""
	for importPath, val := range cusMap {
		importConfig = fmt.Sprintf("%s %s", importConfig, importPath)
		reg := regexp.MustCompile("\"([a-zA-Z]\\w*)\":")
		val = reg.ReplaceAllString(val, `$1:`)
		// schema: "customer"}
		val = strings.ReplaceAll(val, "schema:\"", "schema:")
		val = strings.ReplaceAll(val, "\"}", "}")
		cusConfig = fmt.Sprintf("%s %s,", cusConfig, val)
	}
	config = fmt.Sprintf(`graphqlServers: [%s,%s]`, defaultServer, cusConfig)
	return importConfig, config
}
