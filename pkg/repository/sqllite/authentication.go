package sqllite

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/fire_boom/domain"
	log "github.com/sirupsen/logrus"
)

type authenticationRepository struct {
	db    *sql.DB
	table string
}

func NewAuthenticationRepository(Conn *sql.DB) *authenticationRepository {
	return &authenticationRepository{
		db:    Conn,
		table: "`fb_authentication`",
	}
}

func (a *authenticationRepository) Store(ctx context.Context, f *domain.FbAuthentication) (affect int64, err error) {
	query := fmt.Sprintf("INSERT INTO %s (`name`, `auth_supplier`, `config`, `switch_state`) values(?,?,?,?)", a.table)
	stmt, err := a.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error("authenticationRepository Store PrepareContext err : ", err.Error())
		return
	}
	res, err := stmt.ExecContext(ctx, f.Name, f.AuthSupplier, f.Config, f.SwitchState)
	if err != nil {
		log.Error("authenticationRepository Store ExecContext err : ", err.Error())
		return
	}
	affect, err = res.LastInsertId()
	if err != nil {
		log.Error("authenticationRepository Store res.LastInsertId() err : ", err.Error())
		return
	}
	return
}

func (a *authenticationRepository) Update(ctx context.Context, f *domain.FbAuthentication) (affect int64, err error) {
	query := fmt.Sprintf("UPDATE %s SET `name` = ? , `auth_supplier` = ? ,`config` = ?, `switch_state` = ? WHERE `id` = ?", a.table)
	stmt, err := a.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error("authenticationRepository Update PrepareContext err : ", err.Error())
		return
	}
	res, err := stmt.ExecContext(ctx, f.Name, f.AuthSupplier, f.Config, f.SwitchState, f.ID)
	if err != nil {
		log.Error("authenticationRepository Update ExecContext err : ", err.Error())
		return
	}
	affect, err = res.RowsAffected()
	if err != nil {
		log.Error("authenticationRepository Update RowsAffected err : ", err.Error())
		return
	}

	if affect != rowAffect {
		err = fmt.Errorf("Weird  Behavior. Total Affected: %d", affect)
		return
	}

	return
}

func (a *authenticationRepository) CheckExist(ctx context.Context, f *domain.FbAuthentication) (result domain.FbAuthentication, err error) {
	query := fmt.Sprintf("select `id`,`name`,`auth_supplier`,`switch_state`,`config`,`create_time`,`update_time`,`is_del` from %s where `is_del` != 1 and `id` != ? and `name` = ?", a.table)
	var queryOne []domain.FbAuthentication
	queryOne, err = a.fetch(ctx, query, f.ID, f.Name)
	if err != nil {
		log.Error("authenticationRepository CheckExist fetch err : ", err.Error())
		return
	}
	if len(queryOne) > 0 {
		result = queryOne[0]
	}
	return
}

func (a *authenticationRepository) GetByName(ctx context.Context, name string) (result domain.FbAuthentication, err error) {
	query := fmt.Sprintf("select `id`,`name`,`auth_supplier`,`switch_state`,`config`,`create_time`,`update_time`,`is_del` from %s where `is_del` != 1 and `name` = ?", a.table)
	var queryOne []domain.FbAuthentication
	queryOne, err = a.fetch(ctx, query, name)
	if err != nil {
		log.Error("authenticationRepository GetByName fetch err : ", err.Error())
		return
	}
	if len(queryOne) > 0 {
		result = queryOne[0]
	}
	return
}

func (a *authenticationRepository) Delete(ctx context.Context, id uint) (affect int64, err error) {
	query := fmt.Sprintf("UPDATE %s SET `is_del` = %d where `id` = ? ", a.table, isDel)
	stmt, err := a.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error("authenticationRepository Delete PrepareContext err : ", err.Error())
		return
	}
	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		log.Error("authenticationRepository Delete ExecContext err : ", err.Error())
		return
	}
	affect, err = res.RowsAffected()
	if err != nil {
		log.Error("authenticationRepository Delete RowsAffected err : ", err.Error())
		return
	}

	return
}
func (a *authenticationRepository) FindAuthentication(ctx context.Context) (result []domain.FbAuthentication, err error) {
	query := fmt.Sprintf("SELECT `id`, `name`, `auth_supplier`,`switch_state`, `config`, `create_time`, `update_time`, `is_del` from  %s where `is_del` != %d ", a.table, isDel)
	result, err = a.fetch(ctx, query)
	if err != nil {
		log.Error("authenticationRepository FindAuthentication fetch err : ", err.Error())
		return
	}
	return
}

func (a *authenticationRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.FbAuthentication, err error) {
	rows, err := a.db.QueryContext(ctx, query, args...)
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

	result = make([]domain.FbAuthentication, 0)
	for rows.Next() {
		auth := domain.FbAuthentication{}
		err = rows.Scan(
			&auth.ID,
			&auth.Name,
			&auth.AuthSupplier,
			&auth.SwitchState,
			&auth.Config,
			&auth.CreateTime,
			&auth.UpdateTime,
			&auth.IsDel,
		)

		if err != nil {
			return
		}
		result = append(result, auth)
	}

	return result, nil
}
