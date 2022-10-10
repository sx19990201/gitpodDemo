package sqllite

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/fire_boom/domain"
	log "github.com/sirupsen/logrus"
)

type storageBucketRepository struct {
	db    *sql.DB
	table string
}

func NewStorageBucketRepository(Conn *sql.DB) *storageBucketRepository {
	return &storageBucketRepository{
		db:    Conn,
		table: "`fb_storage_bucket`",
	}
}

func (s *storageBucketRepository) Store(ctx context.Context, f *domain.FbStorageBucket) (result int64, err error) {
	query := fmt.Sprintf("INSERT INTO %s (`name`, `switch`, `config`) values(?,?,?)", s.table)
	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error("storageBucketRepository Store PrepareContext err : ", err.Error())
		return
	}
	res, err := stmt.ExecContext(ctx, f.Name, f.Switch, f.Config)
	if err != nil {
		log.Error("storageBucketRepository Store ExecContext err : ", err.Error())
		return
	}
	result, err = res.LastInsertId()
	if err != nil {
		log.Error("storageBucketRepository Store LastInsertId err : ", err.Error())
		return
	}
	return
}

func (s *storageBucketRepository) Update(ctx context.Context, f *domain.FbStorageBucket) (affect int64, err error) {
	query := fmt.Sprintf("UPDATE %s SET `name` = ? , `switch` = ? ,`config` = ? WHERE `id` = ?", s.table)
	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error("storageBucketRepository Update PrepareContext err : ", err.Error())
		return
	}
	res, err := stmt.ExecContext(ctx, f.Name, f.Switch, f.Config, f.ID)
	if err != nil {
		log.Error("storageBucketRepository Update ExecContext err : ", err.Error())
		return
	}
	affect, err = res.RowsAffected()
	if err != nil {
		log.Error("storageBucketRepository Update LastInsertId err : ", err.Error())
		return
	}

	if affect != rowAffect {
		err = fmt.Errorf("Weird  Behavior. Total Affected: %d", affect)
		return
	}

	return
}

func (s *storageBucketRepository) GetByName(ctx context.Context, name string) (result domain.FbStorageBucket, err error) {
	query := fmt.Sprintf("select `id`,`name`,`switch`,`config`,`create_time`,`update_time`,`is_del` from %s where `is_del` != 1 and `name` = ?", s.table)
	var queryOne []domain.FbStorageBucket
	queryOne, err = s.fetch(ctx, query, name)
	if err != nil {
		log.Error("storageBucketRepository GetByName fetch err : ", err.Error())
		return
	}
	if len(queryOne) > 0 {
		result = queryOne[0]
	}
	return
}

func (s *storageBucketRepository) GetByID(ctx context.Context, id uint) (result domain.FbStorageBucket, err error) {
	query := fmt.Sprintf("select `id`,`name`,`switch`,`config`,`create_time`,`update_time`,`is_del` from %s where `is_del` != 1 and `id` = ?", s.table)
	var queryOne []domain.FbStorageBucket
	queryOne, err = s.fetch(ctx, query, id)
	if err != nil {
		log.Error("storageBucketRepository GetByID fetch err : ", err.Error())
		return
	}
	if len(queryOne) > 0 {
		result = queryOne[0]
	}
	return
}

func (s *storageBucketRepository) CheckExist(ctx context.Context, sb *domain.FbStorageBucket) (result domain.FbStorageBucket, err error) {
	query := fmt.Sprintf("select `id`,`name`,`switch`,`config`,`create_time`,`update_time`,`is_del` from %s where `is_del` != 1 and `id` != ? and `name` = ?", s.table)
	var queryOne []domain.FbStorageBucket
	queryOne, err = s.fetch(ctx, query, sb.ID, sb.Name)
	if err != nil {
		log.Error("storageBucketRepository CheckExist fetch err : ", err.Error())
		return
	}
	if len(queryOne) > 0 {
		result = queryOne[0]
	}
	return
}

func (s *storageBucketRepository) Delete(ctx context.Context, id uint) (affect int64, err error) {
	query := fmt.Sprintf("UPDATE %s SET `is_del` = %d where `id` = ? ", s.table, isDel)
	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error("storageBucketRepository Delete PrepareContext err : ", err.Error())
		return
	}
	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		log.Error("storageBucketRepository Delete ExecContext err : ", err.Error())
		return
	}
	affect, err = res.RowsAffected()
	if err != nil {
		log.Error("storageBucketRepository Delete RowsAffected err : ", err.Error())
		return
	}

	return
}

func (s *storageBucketRepository) FindStorageBucket(ctx context.Context) (result []domain.FbStorageBucket, err error) {
	query := fmt.Sprintf("SELECT `id`, `name`, `switch`, `config`, `create_time`, `update_time`, `is_del` from  %s where `is_del` != %d ", s.table, isDel)
	result, err = s.fetch(ctx, query)
	if err != nil {
		log.Error("storageBucketRepository FindStorageBucket fetch err : ", err.Error())
		return
	}
	return
}

func (s *storageBucketRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.FbStorageBucket, err error) {
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Error("storageBucketRepository fetch QueryContext err : ", err.Error())
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			err = errRow
			return
		}
	}()

	result = make([]domain.FbStorageBucket, 0)
	for rows.Next() {
		dataSource := domain.FbStorageBucket{}
		err = rows.Scan(
			&dataSource.ID,
			&dataSource.Name,
			&dataSource.Switch,
			&dataSource.Config,
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
