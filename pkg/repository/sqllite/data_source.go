package sqllite

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/fire_boom/domain"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"

	log "github.com/sirupsen/logrus"
)

const (
	rowAffect = 1
	isDel     = 1
)

type dataSourceRepository struct {
	db    *sql.DB
	table string
}

func NewDataSourceRepository(Conn *sql.DB) *dataSourceRepository {
	return &dataSourceRepository{
		db:    Conn,
		table: "`fb_data_source`",
	}
}

func (d *dataSourceRepository) GetByName(ctx context.Context, name string) (result domain.FbDataSource, err error) {
	query := fmt.Sprintf("select `id`,`name`,`source_type`,`config`,`switch`,`create_time`,`update_time`,`is_del` from %s where `is_del` != 1 and `name` = ?", d.table)
	var queryOne []domain.FbDataSource
	queryOne, err = d.fetch(ctx, query, name)
	if err != nil {
		log.Error("dataSourceRepository GetByName fetch err : ", err.Error())
		return
	}
	if len(queryOne) > 0 {
		result = queryOne[0]
	}
	return
}

func (d *dataSourceRepository) GetByID(ctx context.Context, id uint) (result domain.FbDataSource, err error) {
	query := fmt.Sprintf("select `id`,`name`,`source_type`,`config`,`switch`,`create_time`,`update_time`,`is_del` from %s where `is_del` != 1 and `id` = ?", d.table)
	var queryOne []domain.FbDataSource
	queryOne, err = d.fetch(ctx, query, id)
	if err != nil {
		log.Error("dataSourceRepository GetByID fetch err : ", err.Error())
		return
	}
	if len(queryOne) > 0 {
		result = queryOne[0]
	}
	return
}

func (d *dataSourceRepository) CheckExist(ctx context.Context, auth *domain.FbDataSource) (result domain.FbDataSource, err error) {
	query := fmt.Sprintf("select `id`,`name`,`source_type`,`config`,`switch`,`create_time`,`update_time`,`is_del` from %s where `is_del` != 1 and `id` != ? and `name` = ?", d.table)
	var queryOne []domain.FbDataSource
	queryOne, err = d.fetch(ctx, query, auth.ID, auth.Name)
	if err != nil {
		log.Error("dataSourceRepository CheckExist fetch err : ", err.Error())
		return
	}
	if len(queryOne) > 0 {
		result = queryOne[0]
	}
	return
}

func (d *dataSourceRepository) Store(ctx context.Context, f *domain.FbDataSource) (result int64, err error) {
	query := fmt.Sprintf("INSERT INTO %s (`name`, `source_type`, `config`,`switch`) values(?,?,?,?)", d.table)
	stmt, err := d.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error("dataSourceRepository Store PrepareContext err : ", err.Error())
		return
	}
	res, err := stmt.ExecContext(ctx, f.Name, f.SourceType, f.Config, f.Switch)
	if err != nil {
		log.Error("dataSourceRepository Store ExecContext err : ", err.Error())
		return
	}
	result, err = res.LastInsertId()
	if err != nil {
		log.Error("dataSourceRepository Store LastInsertId err : ", err.Error())
		return
	}
	return
}

func (d *dataSourceRepository) Update(ctx context.Context, f *domain.FbDataSource) (affect int64, err error) {
	query := fmt.Sprintf("UPDATE %s SET `name` = ? , `source_type` = ? ,`config` = ?,`switch`= ? WHERE `id` = ?", d.table)
	stmt, err := d.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error("dataSourceRepository Update PrepareContext err : ", err.Error())
		return
	}
	res, err := stmt.ExecContext(ctx, f.Name, f.SourceType, f.Config, f.Switch, f.ID)
	if err != nil {
		log.Error("dataSourceRepository Update ExecContext err : ", err.Error())
		return
	}
	affect, err = res.RowsAffected()
	if err != nil {
		log.Error("dataSourceRepository Update LastInsertId err : ", err.Error())
		return
	}

	if affect != rowAffect {
		err = fmt.Errorf("Weird  Behavior. Total Affected: %d", affect)
		return
	}

	return
}

func (d *dataSourceRepository) Delete(ctx context.Context, id uint) (affect int64, err error) {
	query := fmt.Sprintf("UPDATE %s SET `is_del` = %d where `id` = ? ", d.table, isDel)
	stmt, err := d.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error("dataSourceRepository Delete PrepareContext err : ", err.Error())
		return
	}
	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		log.Error("dataSourceRepository Delete ExecContext err : ", err.Error())
		return
	}
	affect, err = res.RowsAffected()
	if err != nil {
		log.Error("dataSourceRepository Delete RowsAffected err : ", err.Error())
		return
	}

	return
}

func (d *dataSourceRepository) FindDataSources(ctx context.Context) (result []domain.FbDataSource, err error) {
	query := fmt.Sprintf("SELECT `id`, `name`, `source_type`, `config`,`switch`, `create_time`, `update_time`, `is_del` from  %s where `is_del` != %d ", d.table, isDel)
	result, err = d.fetch(ctx, query)
	if err != nil {
		log.Error("dataSourceRepository FindDataSources fetch err : ", err.Error())
		return
	}
	return
}

func (d *dataSourceRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.FbDataSource, err error) {
	rows, err := d.db.QueryContext(ctx, query, args...)
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

	result = make([]domain.FbDataSource, 0)
	for rows.Next() {
		dataSource := domain.FbDataSource{}
		err = rows.Scan(
			&dataSource.ID,
			&dataSource.Name,
			&dataSource.SourceType,
			&dataSource.Config,
			&dataSource.Switch,
			&dataSource.CreateTime,
			&dataSource.UpdateTime,
			&dataSource.IsDel,
		)

		if err != nil {
			return
		}
		result = append(result, dataSource)
	}

	return result, nil
}
