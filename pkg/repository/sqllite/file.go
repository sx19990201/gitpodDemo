package sqllite

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/fire_boom/domain"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

type fileRepository struct {
	db    *sql.DB
	table string
}

func NewFileRepository(Conn *sql.DB) *fileRepository {
	return &fileRepository{
		db:    Conn,
		table: "`fb_file`",
	}
}

// FileRead 读取文件所有数据
func (f *fileRepository) FileRead(ctx context.Context, path string) (content []byte, err error) {
	return ioutil.ReadFile(path)
}

// Store 上传文件到 static 目录，meta 信息存储到表
func (f *fileRepository) Store(ctx context.Context, file *domain.File) (result string, err error) {
	query := fmt.Sprintf("INSERT INTO %s (`id`,`name`, `path`) values(?,?,?)", f.table)
	stmt, err := f.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error("fileRepository Store PrepareContext err : ", err.Error())
		return
	}
	_, err = stmt.ExecContext(ctx, file.ID, file.Name, file.Path)
	if err != nil {
		log.Error("fileRepository Store ExecContext err : ", err.Error())
		return
	}
	result = file.ID
	return
}

func (f *fileRepository) Delete(ctx context.Context, fileUUID string) (err error) {
	query := fmt.Sprintf("UPDATE %s SET `is_del` = 1  WHERE `id` = ?", f.table)
	stmt, err := f.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error("fileRepository Update PrepareContext err : ", err.Error())
		return
	}
	_, err = stmt.ExecContext(ctx, fileUUID)
	if err != nil {
		log.Error("fileRepository Delete ExecContext err : ", err.Error())
		return
	}
	return
}

func (f *fileRepository) GetByID(ctx context.Context, fileUUID uint) (result domain.File, err error) {
	query := fmt.Sprintf("select `id`,`code`,`remark`,`create_time`,`update_time`,`is_del` from %s where `is_del` != 1", f.table)
	var queryOne []domain.File
	queryOne, err = f.fetch(ctx, query, fileUUID)
	if err != nil {
		log.Error("roleRepository GetByCode fetch err : ", err.Error())
		return
	}
	if len(queryOne) > 0 {
		result = queryOne[0]
	}
	return
}

func (f *fileRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.File, err error) {
	rows, err := f.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Error("fileRepository fetch QueryContext err : ", err.Error())
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			err = errRow
			return
		}
	}()

	result = make([]domain.File, 0)
	for rows.Next() {
		role := domain.File{}
		err = rows.Scan(
			&role.ID,
			&role.Name,
			&role.Path,
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
