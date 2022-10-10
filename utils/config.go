package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/fire_boom/domain"
	fireboomViper "github.com/fire_boom/pkg/viper"
	log "github.com/sirupsen/logrus"
)

const (
	OFF = false
	ON  = true
)

// 为了避免大量改动, 使用 fireboomViper 替换 github/viper
var viper = fireboomViper.GetFireBoomViper()

func GetOIDCProviderPort() string {
	return viper.GetString("oidcProvider.port")
}

func GetOIDCQuerySchemaPath() string {
	return viper.GetString("dir.oidc.path")
}
func GetOauthLinkerPath() string {
	return viper.GetString("file.oauthLinkerPath")
}

func GetIndexDBIDJsonPath() string {
	return viper.GetString("file.indexDBIDJsonPath")
}

func GetOauthDefaultLinkerPath() string {
	return viper.GetString("file.oauthLinkerDefaultPath")
}
func GetOauthLoginBrandPath() string {
	return viper.GetString("file.oauthLoginBrandPath")
}
func GetOauthLoginConfigPath() string {
	return viper.GetString("file.oauthLoginConfigPath")
}
func GetOauthConfigPath() string {
	return viper.GetString("file.oauthOauthConfigPath")
}

func GetFireBoomPort() string {
	return viper.GetString("server.address")
}

func GetModuleSchemaPath() string {
	return viper.GetString("file.moduleSchemaPath")
}

func GetDefaultOauthDbPath() string {
	return viper.GetString("file.defaultOauthDBPath")
}

func GetPrismaSchemaFilePath() string {
	return fmt.Sprintf("%s/%s", viper.GetString("dir.schema.path"), viper.GetString("file.schemaPath"))
}
func GetPrismaDBPath() string {
	return viper.GetString("dir.schema.dbPath")
}
func GetFireBoomHost() string {
	return viper.GetString("server.host")
}

func GetSdkSrcPath() string {
	return viper.GetString("dir.sdk.src")
}
func GetSdkDstPath() string {
	return viper.GetString("dir.sdk.dst")
}

func GetHost() string {
	return viper.GetString("server.host")
}

func GetWdgStartDir() string {
	return viper.GetString("dir.wdgStartDir")
}

func GetWdgConfigTPL() string {
	return viper.GetString("file.wunderGraphSyncer.configTPL")
}
func GetWdgConfigFile() string {
	return viper.GetString("file.wunderGraphSyncer.configFile")
}
func GetWdgOperationTPL() string {
	return viper.GetString("file.wunderGraphSyncer.operationTPL")
}
func GetWdgOperationFile() string {
	return viper.GetString("file.wunderGraphSyncer.operationFile")
}
func GetWdgServerTPL() string {
	return viper.GetString("file.wunderGraphSyncer.serverTPL")
}
func GetWdgServerFile() string {
	return viper.GetString("file.wunderGraphSyncer.serverFile")
}

func GetTimeOut() int {
	return viper.GetInt("context.timeout")
}

func GetApiPathPrefix() string {
	return viper.GetString("dir.operateAPIPath.prefix")
}

func GetSwitchOff() string {
	return viper.GetString("file.operateAPIPath.switchOff")
}
func GetSwitchOn() string {
	return viper.GetString("file.operateAPIPath.switchOn")
}

func GetApiPathSuffix() string {
	return viper.GetString("file.operateAPIPath.suffix")
}

func GetSettingPath() string {
	return viper.GetString("file.settingPath")
}

func GetWdgGeneratedPath() string {
	return viper.GetString("dir.wunderGraphPath.generated")
}

func GetDriverName() string {
	return viper.GetString("file.database.driverName")
}

func GetDataSourceName() string {
	return viper.GetString("file.database.dataSourceName")
}
func GetOSSUploadPath() string {
	return viper.GetString("dir.ossUploadPath")
}

func GetAuthGlobalHookPathPrefix() string {
	return viper.GetString("dir.Hooks.authGlobalPrefix")
}
func GetNewHookPathPrefix() string {
	return viper.GetString("dir.Hooks.newHookPath")
}

func GetGlobalHookPathPrefix() string {
	return viper.GetString("dir.Hooks.globalPrefix")
}

func GetOperationsHookPathPrefix() string {
	return viper.GetString("dir.Hooks.operationsPrefix")
}

func GetCustomizeHookPathPrefix() string {
	return viper.GetString("dir.Hooks.customizePrefix")
}
func GetHooksSuffix() string {
	return viper.GetString("file.Hooks.suffix")
}

func GetMockPath() string {
	return viper.GetString("dir.Hooks.mockPrefix")
}

func GetOperationSettingPath() string {
	return viper.GetString("dir.operationsSettingPath.operationPrefix")
}

func GetOperationGlobalSettingPath() string {
	return viper.GetString("file.operationsSettingPath.globalPrefix")
}

func GetSwitchState(enable bool) string {
	switch enable {
	case OFF:
		return GetSwitchOff()
	case ON:
		return GetSwitchOn()
	default:
		return GetSwitchOff()
	}
}

func GetOASFilePath() string {
	return viper.GetString("dir.uploadPath.oasPath")
}

func GetDraftPrefixPath() string {
	return viper.GetString("dir.draftPath")
}

func GetJsonSuffix() string {
	return viper.GetString("file.fileType.json")
}
func GetWdgLogPath() string {
	return viper.GetString("file.logPath.wdg")
}
func GetHookLogPath() string {
	return viper.GetString("file.logPath.hook")
}

func GetYamlSuffix() string {
	return viper.GetString("file.fileType.yaml")
}

func GetGlobalConfig() (result domain.GlobalConfig, err error) {
	content, err := ioutil.ReadFile(GetGlobalConfigPath())
	if err != nil {
		log.Error("GetGlobalConfig ioutil.ReadFile err : ", err)
		err = domain.FileReadErr
		return
	}
	err = json.Unmarshal(content, &result)
	if err != nil {
		log.Error("GetGlobalConfig json.Unmarshal err : ", err)
		err = domain.JsonUnMarshalErr
		return
	}
	return
}

func GetGlobalConfigPath() string {
	return viper.GetString("file.globalConfigPath")
}
