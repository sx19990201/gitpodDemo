package sqllite

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/fire_boom/domain"
	log "github.com/sirupsen/logrus"
)

type roleRepository struct {
	db    *sql.DB
	table string
}

func NewRoleRepository(Conn *sql.DB) *roleRepository {
	return &roleRepository{
		db:    Conn,
		table: "`fb_role`",
	}
}

func (r *roleRepository) Store(ctx context.Context, role *domain.FbRole) (affect int64, err error) {
	query := fmt.Sprintf("INSERT INTO %s (`code`, `remark`) values(?,?)", r.table)
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error("roleRepository Store PrepareContext err : ", err.Error())
		return
	}
	res, err := stmt.ExecContext(ctx, role.Code, role.Remark)
	if err != nil {
		log.Error("roleRepository Store ExecContext err : ", err.Error())
		return
	}
	affect, err = res.LastInsertId()
	if err != nil {
		log.Error("roleRepository Store LastInsertId err : ", err.Error())
		return
	}
	return
}

func (r *roleRepository) Update(ctx context.Context, role *domain.FbRole) (affect int64, err error) {
	query := fmt.Sprintf("UPDATE %s SET `code` = ? , `remark` = ? WHERE `id` = ?", r.table)
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error("roleRepository Update PrepareContext err : ", err.Error())
		return
	}
	res, err := stmt.ExecContext(ctx, role.Code, role.Remark, role.ID)
	if err != nil {
		log.Error("roleRepository Update ExecContext err : ", err.Error())
		return
	}
	affect, err = res.RowsAffected()
	if err != nil {
		log.Error("roleRepository Update RowsAffected err : ", err.Error())
		return
	}

	if affect != rowAffect {
		err = fmt.Errorf("Weird  Behavior. Total Affected: %d", affect)
		return
	}

	return
}

func (r *roleRepository) GetByCode(ctx context.Context, name string) (result domain.FbRole, err error) {
	query := fmt.Sprintf("select `id`,`code`,`remark`,`create_time`,`update_time`,`is_del` from %s where `is_del` != 1 and `code` = ?", r.table)
	var queryOne []domain.FbRole
	queryOne, err = r.fetch(ctx, query, name)
	if err != nil {
		log.Error("roleRepository GetByCode fetch err : ", err.Error())
		return
	}
	if len(queryOne) > 0 {
		result = queryOne[0]
	}
	return
}

func (r *roleRepository) CheckExist(ctx context.Context, role *domain.FbRole) (result domain.FbRole, err error) {
	query := fmt.Sprintf("select `id`,`code`,`remark`,`create_time`,`update_time`,`is_del` from %s where `is_del` != 1 and id != ? and `code` = ?", r.table)
	var queryOne []domain.FbRole
	queryOne, err = r.fetch(ctx, query, role.ID, role.Code)
	if err != nil {
		log.Error("roleRepository CheckExist fetch err : ", err.Error())
		return
	}
	if len(queryOne) > 0 {
		result = queryOne[0]
	}
	return
}

func (r *roleRepository) Delete(ctx context.Context, id uint) (affect int64, err error) {
	query := fmt.Sprintf("UPDATE %s SET `is_del` = %d where `id` = ? ", r.table, isDel)
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error("roleRepository Delete PrepareContext err : ", err.Error())
		return
	}
	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		log.Error("roleRepository Delete ExecContext err : ", err.Error())
		return
	}
	affect, err = res.RowsAffected()
	if err != nil {
		log.Error("roleRepository Delete RowsAffected err : ", err.Error())
		return
	}

	return
}

func (r *roleRepository) FindRoles(ctx context.Context) (result []domain.FbRole, err error) {
	query := fmt.Sprintf("SELECT `id`, `code`, `remark`, `create_time`, `update_time`, `is_del` from  %s where `is_del` != %d ", r.table, isDel)
	result, err = r.fetch(ctx, query)
	if err != nil {
		log.Error("roleRepository FindRoles fetch err : ", err.Error())
		return
	}
	return
}

func (r *roleRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.FbRole, err error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Error("roleRepository fetch QueryContext err : ", err.Error())
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			err = errRow
			return
		}
	}()

	result = make([]domain.FbRole, 0)
	for rows.Next() {
		role := domain.FbRole{}
		err = rows.Scan(
			&role.ID,
			&role.Code,
			&role.Remark,
			&role.CreateTime,
			&role.UpdateTime,
			&role.IsDel,
		)

		if err != nil {
			return
		}
		result = append(result, role)
	}

	return result, nil
}
