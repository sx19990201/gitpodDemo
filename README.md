# fire_boom
## 关于 fork_wundergraph sdk
1: fork 之后的地址 https://github.com/javacode123/wundergraph，开发分支 zjl
2: clone fork 之后的 wundergraph 进行开发，集成接口
3: 在 fire_boom 中通过 replace 命令将 wundergraph 替换为本地项目
3: 调试完毕之后将 wundergraph 发送到远端
4: 通过 go get github.com/javacode123/wundergraph@branch_name 获取版本号码
5: 更新 frie_boom 的 replace 命令，使用远端的 wundergraph
「
replace (
	github.com/wundergraph/wundergraph => github.com/javacode123/wundergraph v0.0.0-20220821132915-026aed6f0d49
	github.com/wundergraph/wundergraph/types => github.com/javacode123/wundergraph/types v0.0.0-20220821132915-026aed6f0d49
)
」