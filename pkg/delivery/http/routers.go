package http

import (
	"database/sql"
	"github.com/fire_boom/domain"
	"github.com/labstack/echo/v4"
)

//func InitRouters(e *echo.Echo, db *sql.DB, dsr domain.DataSourceRepository, ws wundergraph.Wdg) {
//	InitDataSourceRouter(e, dsr, ws)
//	InitAuthenticationUseCaseRouter(e, db)
//	InitRoleUseCaseRouter(e, db)
//	InitPrismaUseCaseRouter(e, db)
//	InitDraftHandler(e, dsr)
//	InitFileHandler(e, db)
//	InitStorageBucketRouter(e, db)
//	InitGlobalConfigHandler(e)
//	InitOperateAPIHandler(e)
//	InitScriptHandler(e)
//	InitSettingHandler(e)
//	InitS3UploadClientUseCaseRouter(e, db)
//}

func InitRouters(e *echo.Echo, db *sql.DB, dsr domain.DataSourceRepository) {
	InitDataSourceRouter(e, db, dsr)
	InitAuthenticationUseCaseRouter(e, db)
	InitRoleUseCaseRouter(e, db)
	InitPrismaUseCaseRouter(e, db)
	InitDraftHandler(e, dsr)
	// InitDraftHandler(e, dsr)
	InitFileHandler(e, db)
	InitStorageBucketRouter(e, db)
	InitGlobalConfigHandler(e)
	InitOperateAPIHandler(e, db)
	InitScriptHandler(e)
	InitSettingHandler(e)
	InitS3UploadClientUseCaseRouter(e, db)
	InitEnvRouter(e, db)
	InitLinkerHandler(e)
	InitWdgHandler(e, db)
	InitHomeHandler(e, db)
	InitOauthPrismaUseCaseRouter(e, dsr)
	InitHookHandler(e)
}
