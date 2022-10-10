package wundergraph

import (
	"context"
	"github.com/fire_boom/domain"
	"github.com/fire_boom/domain/mocks"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestReloadConfig(t *testing.T) {
	ctx := context.TODO()
	configTPLPath := "../../static/tpl/config.txt"
	configTSPath := "../../static/tpl/config_test.ts"
	mockDataSourceRepo := new(mocks.DataSourceRepository)
	repository := Repository{
		//dsr: mockDataSourceRepo,
	}
	t.Run("reload successfully", func(t *testing.T) {
		result := []domain.FbDataSource{
			//{
			//	ID:         1,
			//	Name:       "mysqlTest1",
			//	SourceType: 1,
			//	Config:     "{\"apiNamespace\":\"mysqlTest\",\"dbType\":2,\"appendType\":1,\"databaseUrl\":\"\",\"host\":\"8.142.115.204\",\"dbName\":\"gotrue_development\",\"port\":\"3306\",\"userName\":\"root\",\"password\":\"shaoxiong123456\"}\n",
			//	Switch:     domain.SwitchOn,
			//},
			//{
			//	ID:         2,
			//	Name:       "mysqlTest2",
			//	SourceType: 1,
			//	Config:     "{\"apiNamespace\": \"mysqlTest\",\"dbType\": 2,\"appendType\": 2,\"databaseUrl\": \"root:shaoxiong123456@tcp(8.142.115.204:3306)/gotrue_development?parseTime=true&multiStatements=true\",\"host\": \"8.142.115.204\",\"dbName\": \"gotrue_development\",\"port\": \"3306\",\"userName\": \"root\",\"password\": \"shaoxiong123456\"\n}",
			//	Switch:     domain.SwitchOn,
			//},
			{
				ID:         3,
				Name:       "restTest",
				SourceType: 2,
				Config:     "{\"apiNamespace\":\"restTest\",\"header\":[{\"kind\":\"0\",\"key\":\"header1\",\"val\":\"header1\"},{\"key\":\"header2\",\"val\":\"header2\",\"kind\":\"0\"},{\"key\":\"header3\",\"val\":\"header3\",\"kind\":\"0\"},{\"key\":\"header4\",\"kind\":\"2\",\"val\":\"header4\"}],\"statusCodeUnions\":true,\"jwtType\":\"0\",\"secret\":{\"kind\":\"0\",\"val\":\"3212das1d23as1d23as1d32asd\"},\"signingMethod\":\"HS256\"}",
				Switch:     domain.SwitchOn,
			},
			{
				ID:         4,
				Name:       "graphqlTest",
				SourceType: 3,
				Config:     "{\"apiNameSpace\":\"graphqlTest\",\"url\":\"1231231231231322\",\"agreement\":false,\"headers\":[{\"kind\":\"0\",\"key\":\"header1\",\"val\":\"header1\"},{\"key\":\"header2\",\"val\":\"header2\",\"kind\":\"0\"},{\"key\":\"header3\",\"val\":\"header3\",\"kind\":\"0\"},{\"key\":\"header4\",\"val\":\"header4\",\"kind\":\"2\"}],\"internal\":true,\"customFloatScalars\":\"1.1,1.2,1.3\",\"customIntScalars\":\"1,2,3,4\",\"skipRenameRootFields\":\"sfasdf\"}",
				Switch:     domain.SwitchOn,
			},
			/*{
				ID:         1,
				Name:       "1",
				SourceType: 3,
				Config:     "sqlLite",
				Switch:     domain.SwitchOn,
			},
			{
				ID:         1,
				Name:       "1",
				SourceType: 3,
				Config:     "sqlLite1",
				Switch:     domain.SwitchOn,
			},*/
		}

		mockDataSourceRepo.On("FindDataSources", mock.Anything).Return(result, nil).Once()
		wg, err := NewWdg(configTPLPath, configTSPath, configTPLPath, configTSPath, configTPLPath, configTSPath, &repository)
		if err != nil {
			t.Error("unexpect err", err)
		}
		wg.ReloadConfig(ctx)
		//mockDataSourceRepo.On("FindDataSources", mock.Anything).Return(result[0:1], nil)
		//wg.ReloadConfig(ctx)
		//wg.ReloadConfig(ctx)
		//wg.ReloadConfig(ctx)
	})

}
