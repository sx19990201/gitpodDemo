package oidcprovider

import (
	"github.com/fire_boom/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func StartOIDCProvider() {
	// TODO 拉取最新的db结构
	//err := cli.Run([]string{"generate", "--schema", path}, true)
	//if err != nil {
	//	fmt.Println(err)
	//}

	// 注册路由
	RegisterHandlers()
	// 初始化数据库查询schema
	utils.InitQuerySchema(utils.GetOIDCQuerySchemaPath())
	port := utils.GetOIDCProviderPort()

	//_ = exec.Command("open", "http://localhost:"+port).Run()
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
