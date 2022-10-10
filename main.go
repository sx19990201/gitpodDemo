package main

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/fire_boom/pkg/oidcprovider"
	"github.com/fire_boom/pkg/usecase"
	"github.com/fire_boom/pkg/wundergraph"
	"github.com/labstack/echo/v4"
	"github.com/prisma/prisma-client-go/engine"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/pprof"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	_ "github.com/fire_boom/app/docs"
	"github.com/fire_boom/domain"
	fireboomHttp "github.com/fire_boom/pkg/delivery/http"
	"github.com/fire_boom/pkg/delivery/http/middleware"
	"github.com/fire_boom/pkg/repository/sqllite"
	fireboomViper "github.com/fire_boom/pkg/viper"

	"github.com/fire_boom/utils"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "net/http/pprof"
)

// @title  FireBoom-server Swagger API
// @version 1.0
// @description FireBoom-server

// @host localhost:9123
// @produce json

var GlobalStartTime string

func main() {

	////初始化log实例
	//timeLog := fireboomHttp.GetWdgLog()
	//// 创建一个hook，将请求id存进去做该请求日志的唯一标识
	//hook := fireboomHttp.WdgLogHook{}
	//fireboomHttp.WdgLog.AddHook(hook)
	//timeLog.Info("start")

	// 新建一个 viper, 防止与 wundergraph pkg 的 viper 冲突
	v := fireboomViper.GetFireBoomViper()
	err := initConfig(v)

	// ctrl+c
	closeHandler()

	// 初始化目录
	CheckConfigDir(filepath.Join(utils.GetRootPath(), "config.yaml"))
	CheckConfigFile()
	GlobalStartTime = time.Now().Format("2006-01-02 15:04:05")

	checkDependency()
	if err != nil {
		panic(fmt.Sprintf("init config file err=%+v", err))
	}

	initLogs()

	//go func() {
	//	initGoTrue()
	//}()

	db, err := initDB(v)
	if err != nil {
		panic(fmt.Sprintf("init db meet err=%+v", err))
	}

	// 初始化订阅
	//startMessage := pubsub.GetTopicMessage(pubsub.GetPubSub(), pubsub.TopicName)
	//go Process(pubsub.GetMessageChan(startMessage))

	fireboomHttp.StartCount++
	err = wundergraph.ReloadWdgFile(db)
	// 重试一次
	if err != nil {
		wundergraph.ReloadWdgFile(db)
	}

	fireboomHttp.WdgStatus = fireboomHttp.Stop
	// 启动
	go fireboomHttp.StartWdg()

	//fireboomHttp.PushMsg()
	fireboomHttp.StatusChan <- fireboomHttp.WdgStatus

	fireboomHttp.PushLogMsg()
	fireboomHttp.PushStatusMsg()

	dataSourceRepo := sqllite.NewDataSourceRepository(db)

	dbID := utils.GetIndexDBID()
	dbUC := usecase.NewDataSourceUseCase(dataSourceRepo, time.Second)
	go func() {
		schemaContent := ""
		// 如果是默认数据源
		if dbID == -1 {
			schemaContent = utils.GetDefaultDBSchema()
		} else {
			// 组织该数据源的schema文件
			content, err := dbUC.GetPrismaSchema(context.Background(), uint(dbID))
			if err != nil {
				return
			}
			schemaContent = content
		}
		engine.InitQueryEngine(schemaContent)
		//if strings.Contains(err.Error(), "introspect error: The introspected database was empty") {
		//	engine.InitQueryEngine(schemaContent)
		//}
	}()
	// 启动oidcProvider
	go oidcprovider.StartOIDCProvider()

	initHttpSvr(db, v, dataSourceRepo)

}

func closeHandler() {
	quit := make(chan os.Signal, 2)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		for {
			signal := <-quit
			log.Info("Received interrupt signal", signal.String())
			os.Exit(-1)
		}
	}()
}

func initConfig(v *viper.Viper) error {
	path := filepath.Join(utils.GetRootPath(), "config.yaml")
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		return err
	}

	if v.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}
	return nil
}

func initLogs() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func checkDependency() {
	c := exec.Command("node", "-v")
	_, err := c.CombinedOutput()
	if err != nil {
		panic("please install node!!!")
	}
}

func initGoTrue() {
	// windows 系统可能要把 bash 改为 cmd?
	cmdType := "bash"
	if runtime.GOOS == "windows" {
		cmdType = "cmd"
	}
	c := exec.Command(cmdType, "-g", "cd gotrue-master && go run main.go")
	stdout, err := c.StdoutPipe()
	if err != nil {
		panic(fmt.Sprintf("init goTrue stdoutPipe meet err=%+v", err))
	}

	go func() {
		reader := bufio.NewReader(stdout)
		for {
			readString, err := reader.ReadString('\n')
			if err != nil || err == io.EOF {
				break
			}
			log.Println(readString)
		}
	}()
	if err := c.Run(); err != nil {
		panic(fmt.Sprintf("init goTrue run cmd meet err=%+v", err))
	}
}

func initDB(v *viper.Viper) (*sql.DB, error) {
	driverName := utils.GetDriverName()
	dataSourceName := utils.GetDataSourceName()
	db, err := sql.Open(driverName, dataSourceName)
	return db, err
}

func initHttpSvr(db *sql.DB, v *viper.Viper, dsr domain.DataSourceRepository) {
	e := echo.New()

	middL := middleware.InitMiddleware(db)
	e.Use(middL.CORS)
	e.Use(echoMiddleware.Logger())
	e.Use(middL.UpdateWunderGraph)
	e.Use(middL.RequestID)

	fireboomHttp.InitRouters(e, db, dsr)
	log.WithFields(log.Fields{
		"animal": "walrus",
	}).Info("A walrus appears")
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/api/v1/setting/getTime", func(c echo.Context) error {
		return c.JSON(200, fireboomHttp.SuccessWriteResult(GlobalStartTime))
	})

	oasPath := filepath.ToSlash(fmt.Sprintf("%s/%s", utils.GetRootPath(), utils.GetOASFilePath()))
	e.Static(fmt.Sprintf("/%s", utils.GetOASFilePath()), oasPath)

	// 静态资源
	path := filepath.ToSlash(fmt.Sprintf("%s/static/front", utils.GetRootPath()))
	e.Static("/", path)

	e.GET("/app/main/graphql", func(c echo.Context) error {
		resp, err := fireboomHttp.DoProxy(c.Request(), "http://localhost:9123")
		if err != nil {
			return c.JSON(http.StatusOK, "request fail  ")
		}
		return c.HTML(http.StatusOK, resp)
		//return c.Redirect(http.StatusMovedPermanently, "http://8.142.115.204:9991/app/main/graphql")
	})

	e.POST("/app/main/graphql", func(c echo.Context) error {
		resp, err := fireboomHttp.DoProxy(c.Request(), "http://localhost:9123")
		if err != nil {
			return c.JSON(http.StatusOK, "request fail  ")
		}
		return c.String(http.StatusOK, resp)
	})

	e.GET("/debug/pprof", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
	e.GET("/debug/pprof/allocs", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
	e.GET("/debug/pprof/block", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
	e.GET("/debug/pprof/goroutine", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
	e.GET("/debug/pprof/heap", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
	e.GET("/debug/pprof/mutex", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
	e.GET("/debug/pprof/cmdline", echo.WrapHandler(http.HandlerFunc(pprof.Cmdline)))
	e.GET("/debug/pprof/profile", echo.WrapHandler(http.HandlerFunc(pprof.Profile)))
	e.GET("/debug/pprof/symbol", echo.WrapHandler(http.HandlerFunc(pprof.Symbol)))
	e.GET("/debug/pprof/trace", echo.WrapHandler(http.HandlerFunc(pprof.Trace)))

	address := fmt.Sprintf(":%s", utils.GetFireBoomPort())
	e.Start(address)
	//e.StartAutoTLS(address)
	//e.StartTLS(address, "cert.pem", "key.pem")
	//e.StartTLS(address, "cert.pem", "key.pem")
}

type DirPathStruct struct {
	PathArr []string
}

func CheckConfigDir(path string) {
	yamlDir := struct {
		Dir map[string]interface{}
	}{}
	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}
	if err = yaml.Unmarshal(content, &yamlDir); err != nil {
		fmt.Println(err)
	}
	dirPathArr := DirPathStruct{
		PathArr: make([]string, 0),
	}
	GetPaths(yamlDir.Dir, &dirPathArr)
	for _, row := range dirPathArr.PathArr {
		// 不存在则创建
		if !utils.FileExist(row) {
			err = os.MkdirAll(row, os.ModePerm)
		}
	}
}

func CheckConfigFile() {
	if !utils.FileExist(utils.GetGlobalConfigPath()) {
		globalConfig := domain.GlobalConfig{
			AuthRedirectURL: make([]string, 0),
			Hooks: domain.HooksConfiguration{
				Authentication: domain.Authentication{},
				RestApi:        domain.RestApi{},
			},
			ConfigureWunderGraphApplication: domain.WunderGraphConfigApplicationConfig{
				Security: &domain.SecurityConfig{
					EnableGraphQLEndpoint: true,
					AllowedHosts:          make([]string, 0),
				},
				Cors: &domain.CorsConfiguration{
					AllowedOrigins: make([]string, 0),
					AllowedMethods: make([]string, 0),
					AllowedHeaders: make([]string, 0),
					ExposedHeaders: make([]string, 0),
				},
			},
		}
		globalConfigByte, _ := json.Marshal(globalConfig)
		utils.WriteFile(utils.GetGlobalConfigPath(), string(globalConfigByte))
	}
	if !utils.FileExist(utils.GetSettingPath()) {
		setting := domain.Setting{
			System: domain.SystemConfig{
				// 默认开启
				ApiPort:          "9991",
				MiddlewarePort:   "9992",
				ForcedJumpSwitch: true,
				DevSwitch:        true,
				DebugSwitch:      true,
			},
			Version: domain.VersionConfig{},
			Environment: domain.EnvironmentConfig{
				EnvironmentList: make([]domain.EnvironmentDetail, 0),
			},
		}
		settingConfigByte, _ := json.Marshal(setting)
		utils.WriteFile(utils.GetSettingPath(), string(settingConfigByte))
	}
	if !utils.FileExist(utils.GetOperationGlobalSettingPath()) {
		operationGlobalSetting := domain.OperationSetting{}
		operationGlobalSettingByte, _ := json.Marshal(operationGlobalSetting)
		utils.WriteFile(utils.GetOperationGlobalSettingPath(), string(operationGlobalSettingByte))
	}
}

func GetPaths(maps map[string]interface{}, pathArr *DirPathStruct) {
	for _, rowVal := range maps {
		str, isString := rowVal.(string)
		rowMap, isMap := rowVal.(map[string]interface{})
		if isString {
			pathArr.PathArr = append(pathArr.PathArr, str)
		} else if isMap {
			GetPaths(rowMap, pathArr)
		}
	}
}
