package sqllite

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/fire_boom/domain"
	log "github.com/sirupsen/logrus"
)

const (
	likeQuery = "%"
)

type envRepository struct {
	db    *sql.DB
	table string
}

func NewEnvRepository(Conn *sql.DB) *envRepository {
	return &envRepository{
		db:    Conn,
		table: "`fb_env`",
	}
}

func (e *envRepository) Store(ctx context.Context, f *domain.FbEnv) (result int64, err error) {
	query := fmt.Sprintf("INSERT INTO %s (`key`, `dev_env`, `pro_env`,`env_type`) values(?,?,?,?)", e.table)
	stmt, err := e.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error("envRepository Store PrepareContext err : ", err.Error())
		return
	}
	res, err := stmt.ExecContext(ctx, f.Key, f.DevEnv, f.ProEnv, f.EnvType)
	if err != nil {
		log.Error("envRepository Store ExecContext err : ", err.Error())
		return
	}
	result, err = res.LastInsertId()
	if err != nil {
		log.Error("envRepository Store LastInsertId err : ", err.Error())
		return
	}
	return
}

func (e *envRepository) Update(ctx context.Context, f *domain.FbEnv) (result int64, err error) {
	query := fmt.Sprintf("UPDATE %s SET `key` = ? , `dev_env` = ? ,`pro_env` = ? WHERE `id` = ?", e.table)
	stmt, err := e.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error("envRepository Update PrepareContext err : ", err.Error())
		return
	}
	res, err := stmt.ExecContext(ctx, f.Key, f.DevEnv, f.ProEnv, f.ID)
	if err != nil {
		log.Error("envRepository Update ExecContext err : ", err.Error())
		return
	}
	result, err = res.RowsAffected()
	if err != nil {
		log.Error("envRepository Update LastInsertId err : ", err.Error())
		return
	}

	if result != rowAffect {
		err = fmt.Errorf("Weird  Behavior. Total Affected: %d", result)
		return
	}

	return
}

func (e *envRepository) GetByKey(ctx context.Context, key string) (result domain.FbEnv, err error) {
	query := fmt.Sprintf("SELECT `id`, `key`, `dev_env`, `pro_env`,`env_type`, `create_time`, `update_time`, `is_del` from  %s where `is_del` != %d and key = ?", e.table, isDel)
	envs, err := e.fetch(ctx, query, key)
	if err != nil {
		log.Error("envRepository GetByName fetch err : ", err.Error())
		return
	}
	if len(envs) > 0 {
		result = envs[0]
	}
	return
}

func (e *envRepository) Exist(ctx context.Context, id int64, key string) (result domain.FbEnv, err error) {
	query := fmt.Sprintf("SELECT `id`, `key`, `dev_env`, `pro_env`,`env_type`, `create_time`, `update_time`, `is_del` from  %s where `is_del` != %d and key = ? and id != ?", e.table, isDel)
	envs, err := e.fetch(ctx, query, key, id)
	if err != nil {
		log.Error("envRepository GetByName fetch err : ", err.Error())
		return
	}
	if len(envs) > 0 {
		result = envs[0]
	}
	return
}

func (e *envRepository) Delete(ctx context.Context, id uint) (result int64, err error) {
	query := fmt.Sprintf("UPDATE %s SET `is_del` = %d where `id` = ? ", e.table, isDel)
	stmt, err := e.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error("envRepository Delete PrepareContext err : ", err.Error())
		return
	}
	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		log.Error("envRepository Delete ExecContext err : ", err.Error())
		return
	}
	result, err = res.RowsAffected()
	if err != nil {
		log.Error("envRepository Delete RowsAffected err : ", err.Error())
		return
	}

	return
}

func (e *envRepository) FindEnvs(ctx context.Context, key string) (result []domain.FbEnv, err error) {
	query := fmt.Sprintf("SELECT `id`, `key`, `dev_env`, `pro_env`,`env_type`, `create_time`, `update_time`, `is_del` from  %s where `is_del` != %d and `key` like '%s%s%s'", e.table, isDel, likeQuery, key, likeQuery)
	result, err = e.fetch(ctx, query)
	if err != nil {
		log.Error("envRepository FindEnvs fetch err : ", err.Error())
		return
	}
	return
}

func (e *envRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.FbEnv, err error) {
	rows, err := e.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			err = errRow
			return
		}
	}()

	result = make([]domain.FbEnv, 0)
	for rows.Next() {
		env := domain.FbEnv{}
		err = rows.Scan(
			&env.ID,
			&env.Key,
			&env.DevEnv,
			&env.ProEnv,
			&env.EnvType,
			&env.CreateTime,
			&env.UpdateTime,
			&env.IsDel,
		)

		if err != nil {
			return
		}
		result = append(result, env)
	}

	return result, nil
}
