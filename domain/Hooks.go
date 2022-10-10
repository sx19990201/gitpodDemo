package domain

var opeartionHooksArr = []string{
	"preResolve",
	"postResolve",
	"customResolve",
	"mutatingPostResolve",
	"mutatingPreResolve",
}
var globalHook = []string{"onRequest", "onResponse"}

var authGlobalHook = []string{"postAuthentication", "mutatingPostAuthentication"}

const (
	OperationHook = iota + 1
	AuthGlobalHook
	GlobalHook
)

type Hooks struct {
	FileName   string `json:"fileName"`
	HookName   string `json:"hookName"`
	HookSwitch bool   `json:"hookSwitch"`
	Content    string `json:"content"`
}

type HooksUseCase interface {
	FindHooksByPath(hookPath, methodPath string, hookType int64) ([]Hooks, error)
	UpdateHooksByPath(path string, hooks Hooks) error
	OperationHooksReName(oldName, newName string) error
	RemoveOperationHooks(operationPath string) error
}

func GetHookArr(hookType int64) []string {
	switch hookType {
	case OperationHook:
		return opeartionHooksArr
	case AuthGlobalHook:
		return authGlobalHook
	case GlobalHook:
		return globalHook
	}
	return nil
}
