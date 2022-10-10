/*
generate or update wundergraph.config.ts
*/
package wundergraph

import (
	"context"
	"database/sql"
	"errors"
	"github.com/fire_boom/domain"
	"github.com/fire_boom/pkg/repository/sqllite"
	"github.com/fire_boom/utils"
	log "github.com/sirupsen/logrus"
	"os"
	"sync"
	"text/template"
)

const (
	queries       = "queries"
	mutations     = "mutations"
	subscriptions = "subscriptions"
)

var OperationTypeMap = map[string]string{"query": queries, "mutation": mutations, "subscription": subscriptions}

var wdgMap sync.Map

type item struct {
	tpl  *template.Template
	fh   *os.File
	lock sync.Mutex
}

type wdg struct {
	configItem   *item
	opertionItem *item
	serverItem   *item
	repository   *Repository
}

type Repository struct {
	dsr domain.DataSourceRepository
	ar  domain.AuthenticationRepository
	rr  domain.RoleRepository
	sbr domain.StorageBucketRepository
	er  domain.EnvRepository
	or  domain.OperationsRepository
}

func NewWdg(configTPLPath, configTSPath, operationTPLPath, opertionTSPath, serverTPLPath, serverTSPath string,
	repository *Repository) (Wdg, error) {
	ct, err := template.ParseFiles(configTPLPath)
	if err != nil {
		return nil, err
	}
	cf, err := os.OpenFile(configTSPath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	ot, err := template.ParseFiles(operationTPLPath)
	if err != nil {
		return nil, err
	}
	of, err := os.OpenFile(opertionTSPath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	st, err := template.ParseFiles(serverTPLPath)
	if err != nil {
		return nil, err
	}
	sf, err := os.OpenFile(serverTSPath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return &wdg{
		configItem:   &item{tpl: ct, fh: cf},
		opertionItem: &item{tpl: ot, fh: of},
		serverItem:   &item{tpl: st, fh: sf},
		repository:   repository,
	}, nil
}

func (w *wdg) ReloadConfig(ctx context.Context) error {
	ci := w.configItem
	ci.lock.Lock()
	defer ci.lock.Unlock()

	configData, err := w.buildConfigData(ctx)
	if err != nil {
		return err
	}

	err = ci.fh.Truncate(0)
	if err != nil {
		return err
	}

	_, err = ci.fh.Seek(0, 0)
	if err != nil {
		return err
	}

	return ci.tpl.Execute(ci.fh, configData)
}

func (w *wdg) ReloadOpertionTs(ctx context.Context) error {
	ci := w.opertionItem
	ci.lock.Lock()
	defer ci.lock.Unlock()

	operationData, err := w.buildOperationData(ctx)
	if err != nil {
		return err
	}

	err = ci.fh.Truncate(0)
	if err != nil {
		return err
	}

	_, err = ci.fh.Seek(0, 0)
	if err != nil {
		return err
	}

	return ci.tpl.Execute(ci.fh, operationData)
}

func (w *wdg) ReloadServerTs(ctx context.Context) error {
	si := w.serverItem
	si.lock.Lock()
	defer si.lock.Unlock()

	serverData, err := w.buildServerData(ctx)
	if err != nil {
		return err
	}

	err = si.fh.Truncate(0)
	if err != nil {
		return err
	}

	_, err = si.fh.Seek(0, 0)
	if err != nil {
		return err
	}

	return si.tpl.Execute(si.fh, serverData)
}

func GetWdg(db *sql.DB) (result Wdg, err error) {
	if val, ok := wdgMap.Load("wdg"); ok {
		result, ok = val.(*wdg)
		if !ok {
			return nil, errors.New("get wdg fail")
		}
		return
	}
	repository := Repository{
		dsr: sqllite.NewDataSourceRepository(db),
		ar:  sqllite.NewAuthenticationRepository(db),
		rr:  sqllite.NewRoleRepository(db),
		sbr: sqllite.NewStorageBucketRepository(db),
		er:  sqllite.NewEnvRepository(db),
		or:  sqllite.NewOperationsRepository(db),
	}

	wdg, err := NewWdg(
		utils.GetWdgConfigTPL(), utils.GetWdgConfigFile(),
		utils.GetWdgOperationTPL(), utils.GetWdgOperationFile(),
		utils.GetWdgServerTPL(), utils.GetWdgServerFile(),
		&repository,
	)
	if err != nil {
		return nil, errors.New("create wdg fail,err : " + err.Error())
	}
	wdgMap.Store("wdg", wdg)
	return wdg, nil
}

func ReloadWdgFile(conn *sql.DB) (err error) {
	wdgRoload, err := GetWdg(conn)
	if err != nil {
		return errors.New("reload Wdg File fail")
	}
	// todo 如果未知原因会panic,先捕获
	defer func() {
		if err := recover(); err != nil {
			err = errors.New("reload Wdg File fail")
		}
	}()
	ctx := context.Background()
	wdgRoload.ReloadConfig(ctx)
	wdgRoload.ReloadOpertionTs(ctx)
	wdgRoload.ReloadServerTs(ctx)
	log.Info("reload config success!")
	return
}
