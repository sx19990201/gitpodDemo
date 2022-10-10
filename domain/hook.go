package domain

type HookStruct struct {
	Path       string      `json:"path"` // 钩子的路径 钩子分四种operation、auth、customize、global 以这4个为一级目录，二级目录为钩子的类型(operation的则以operation的名称为二级目录，钩子类型为三级目录)
	Depends    []Depend    `json:"depend"`
	Script     string      `json:"script"`
	HookSwitch bool        `json:"switch"`
	ScriptType string      `json:"scriptType"`
	Input      interface{} `json:"input"`
	HookType   string      `json:"type"`
}

type Depend struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

var HookDemo = `{
  "path": "operation/FindAll/onRequest.ts",
  "switch": true,
  "depend": [{
   "name": "@type/node",
   "version": "0.1"
  }],
  "script": "console.log(111)",
  "scriptType": "typescript",
  "input": {},
  "type": "preResolve"
}`
