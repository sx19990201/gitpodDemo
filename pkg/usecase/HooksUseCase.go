package usecase

import (
	"fmt"
	"github.com/fire_boom/domain"
	"github.com/fire_boom/utils"
	"github.com/labstack/gommon/log"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type HooksUseCase struct {
	contextTimeout time.Duration
}

func NewHooksUseCase(timeout time.Duration) *HooksUseCase {
	return &HooksUseCase{
		contextTimeout: timeout,
	}
}

func (h *HooksUseCase) FindHooksByPath(hookPath, methodPath string, hookType int64) (result []domain.Hooks, err error) {
	// 获取hooks目录下所有文件
	hookFilePaths := utils.GetDicAllChildFilePath(hookPath)
	hooksArr := domain.GetHookArr(hookType)
	// 按hooks类型遍历
	for index, hookName := range hooksArr {
		// 初始化返回值，默认空字符串
		result = append(result, domain.Hooks{
			HookName:   hookName,
			HookSwitch: utils.OFF, // 默认开关为false，关
		})
		hookFileName := hookName
		// 如果方法路径为空，则说明是全局hook，反之则是operationHook
		if methodPath != "" {
			// operationHook需要加上对应的operation前缀
			// 拼接路径 user/xxxx/xxxx_hooksName
			hookFileName = fmt.Sprintf("%s_%s%s", methodPath, hookName, utils.GetHooksSuffix())
		}
		// 遍历目录下文件
		for _, hookFilePath := range hookFilePaths {
			hookFilePath = filepath.ToSlash(hookFilePath)
			hookFileName = filepath.ToSlash(hookFileName)
			// 如果该文件名包含
			if strings.Contains(hookFilePath, hookFileName) {
				// 拼接路径
				content, err := ioutil.ReadFile(hookFilePath)
				if err != nil {
					log.Error("读取内容失败")
				}
				result[index].Content = string(content)
				if !strings.Contains(hookFilePath, utils.GetSwitchOff()) {
					result[index].HookSwitch = utils.ON
				}
			}
		}
	}
	return
}

func (h *HooksUseCase) UpdateHooksByPath(path string, hooks domain.Hooks) (err error) {
	// 获取hook文件列表
	hookFilePaths := utils.GetDicAllChildFilePath(path)
	// hook文件名,全局hook的文件名就是hook的名字
	thisHookName := hooks.HookName
	// 如果文件名不为空，则说明是operationHook，需要做拼接
	if hooks.FileName != "" {
		// 截取文件名，获得hook文件名
		thisHookName = fmt.Sprintf("%s_%s", strings.Split(hooks.FileName, ".")[0], hooks.HookName)
	}

	// 默认要修改的hook并不存在
	updateHookPath := ""
	// 判断该operation的hook类型的文件是否存在
	for _, hookFilePath := range hookFilePaths {
		hookFilePath = filepath.ToSlash(hookFilePath)
		thisHookName = filepath.ToSlash(thisHookName)
		// 如果存在,则hook文件名改为原文件名
		if strings.Contains(hookFilePath, thisHookName) {
			updateHookPath = hookFilePath
			break
		}
	}
	// 如果不存在,生成文件名，并写入退出
	if updateHookPath == "" {
		if strings.Split(thisHookName, "")[0] != "/" {
			thisHookName = fmt.Sprintf("/%s", thisHookName)
		}
		updateHookPath = fmt.Sprintf("%s%s%s%s", path, thisHookName, utils.GetHooksSuffix(), utils.GetSwitchState(hooks.HookSwitch))
		err = utils.WriteFile(updateHookPath, hooks.Content)
		if err != nil {
			return
		}
		return
	}

	// 如果存在,则根据开关，重新生成名称，并重命名+更改内容
	// 将开关去掉，方便赋值开关，因为开启会加上off直接去掉就行，开启则是空后缀不需要做任何处理
	oldPath := updateHookPath
	newPath := strings.ReplaceAll(updateHookPath, utils.GetSwitchOff(), "")
	newPath = fmt.Sprintf("%s%s", newPath, utils.GetSwitchState(hooks.HookSwitch))
	err = os.Rename(oldPath, newPath)
	if err != nil {
		return
	}
	// 写入文件
	err = utils.WriteFile(newPath, hooks.Content)
	if err != nil {
		return
	}

	return
}

// OperationHooksReName TODO 重命名hooks接口好像只会在operations重命名时调用，现在只处理operation的hooks
func (h *HooksUseCase) OperationHooksReName(oldName, newName string) (err error) {

	// 获取hook文件列表
	hookFilePaths := utils.GetDicAllChildFilePath(utils.GetOperationsHookPathPrefix())
	hooksArr := domain.GetHookArr(domain.OperationHook)

	replaceMap := make(map[string]string, 0)
	// 按hooks类型遍历
	for _, hookName := range hooksArr {
		// 如果是operation，则生成拼接对应的名字
		oldReplaceName := fmt.Sprintf("%s%s_%s%s", utils.GetOperationsHookPathPrefix(), oldName, hookName, utils.GetHooksSuffix())
		newReplaceName := fmt.Sprintf("%s%s_%s%s", utils.GetOperationsHookPathPrefix(), newName, hookName, utils.GetHooksSuffix())
		replaceMap[newReplaceName] = oldReplaceName
	}

	for replaceNewName, replaceOldName := range replaceMap {
		onOldHookPath := filepath.ToSlash(replaceOldName)
		offOldHookPath := filepath.ToSlash(fmt.Sprintf("%s%s", replaceOldName, utils.GetSwitchState(utils.OFF)))
		for _, row := range hookFilePaths {
			row = filepath.ToSlash(row)
			if !(row == onOldHookPath || row == offOldHookPath) {
				continue
			}
			if row == offOldHookPath {
				replaceNewName = fmt.Sprintf("%s%s", replaceNewName, utils.GetSwitchState(false))
			}
			err = os.Rename(row, replaceNewName)
			if err != nil {
				log.Error("OperationHooksReName fail ,err : ", err)
			}
		}
	}

	return
}

// RemoveOperationHooks TODO 删除operations hooks
func (h *HooksUseCase) RemoveOperationHooks(operationPath string) (err error) {

	// 获取hook文件列表
	hookFilePaths := utils.GetDicAllChildFilePath(utils.GetOperationsHookPathPrefix())
	hooksArr := domain.GetHookArr(domain.OperationHook)

	hookArr := make([]string, 0)
	// 按hooks类型遍历
	for _, hookName := range hooksArr {
		oeprHookName := fmt.Sprintf("%s%s_%s%s", utils.GetOperationsHookPathPrefix(), operationPath, hookName, utils.GetHooksSuffix())
		hookArr = append(hookArr, oeprHookName)
	}

	for _, path := range hookArr {
		onPath := filepath.ToSlash(path)
		offPath := filepath.ToSlash(fmt.Sprintf("%s%s", path, utils.GetSwitchState(utils.OFF)))
		for _, row := range hookFilePaths {
			row = filepath.ToSlash(row)
			if !(row == offPath || row == onPath) {
				continue
			}
			os.Remove(row)
		}
	}

	return
}
