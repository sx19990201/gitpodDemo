package sqllite

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/fire_boom/domain"
	"github.com/fire_boom/utils"
	log "github.com/sirupsen/logrus"
	"reflect"
	"strconv"
	"strings"
)

type operationsRepository struct {
	db    *sql.DB
	table string
}

func NewOperationsRepository(Conn *sql.DB) *operationsRepository {
	return &operationsRepository{
		db:    Conn,
		table: "`fb_operations`",
	}
}

func (o *operationsRepository) Store(ctx context.Context, f *domain.FbOperations) (result int64, err error) {
	query := fmt.Sprintf("INSERT INTO %s (`method` ,`operation_type` ,`is_public` ,`remark` ,`legal` ,`path` ,`enable`) values(? ,? ,? ,? ,? ,? ,?)", o.table)
	stmt, err := o.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error("operationsRepository Store PrepareContext err : ", err.Error())
		return
	}
	res, err := stmt.ExecContext(ctx, f.Method, f.OperationType, f.IsPublic, f.Remark, f.Legal, f.Path, f.Enable)
	if err != nil {
		log.Error("operationsRepository Store ExecContext err : ", err.Error())
		return
	}
	result, err = res.LastInsertId()
	if err != nil {
		log.Error("operationsRepository Store LastInsertId err : ", err.Error())
		return
	}
	return
}

func (o *operationsRepository) ReName(ctx context.Context, oldName, newName string, enable int64) (err error) {
	query := fmt.Sprintf("UPDATE %s SET `path` = ?,`enable` = ?  WHERE `path` = ?", o.table)
	stmt, err := o.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error("operationsRepository ReName PrepareContext err : ", err.Error())
		return
	}
	res, err := stmt.ExecContext(ctx, newName, enable, oldName)
	if err != nil {
		log.Error("operationsRepository ReName ExecContext err : ", err.Error())
		return
	}
	_, err = res.RowsAffected()
	if err != nil {
		log.Error("operationsRepository Update LastInsertId err : ", err.Error())
		return
	}

	return
}

func (o *operationsRepository) Update(ctx context.Context, f *domain.FbOperations) (result int64, err error) {
	//query := fmt.Sprintf("UPDATE %s SET `method` = ? ,`operation_type` = ? ,`is_public` = ? ,`remark` = ? ,`legal` = ? ,`enable` = ? ,`content` = ? WHERE `id` = ?", o.table)
	query := GetUpdateSql(*f, o.table)
	stmt, err := o.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error("operationsRepository Update PrepareContext err : ", err.Error())
		return
	}
	//res, err := stmt.ExecContext(ctx, f.Method, f.OperationType, f.IsPublic, f.Remark, f.Legal, f.Enable, f.Content, f.ID)
	res, err := stmt.ExecContext(ctx)
	if err != nil {
		log.Error("operationsRepository Update ExecContext err : ", err.Error())
		return
	}
	result, err = res.RowsAffected()
	if err != nil {
		log.Error("operationsRepository Update LastInsertId err : ", err.Error())
		return
	}

	if result != rowAffect {
		err = fmt.Errorf("Weird  Behavior. Total Affected: %d", result)
		return
	}

	return
}

func (o *operationsRepository) FindByPath(ctx context.Context, path string) (result []domain.FbOperations, err error) {
	query := fmt.Sprintf("SELECT `id` ,`method` ,`operation_type` ,`is_public` ,`remark` ,`legal` ,`path` ,`content` ,`enable` ,`create_time` ,`update_time` ,`is_del` from  %s where `is_del` != %d and `path` like '%s%s'", o.table, isDel, path, likeQuery)
	result, err = o.fetch(ctx, query)
	if err != nil {
		log.Error("operationsRepository GetByPath fetch err : ", err.Error())
		return
	}

	return
}

func (o *operationsRepository) GetByPath(ctx context.Context, path string) (result domain.FbOperations, err error) {
	query := fmt.Sprintf("SELECT `id` ,`method` ,`operation_type` ,`is_public` ,`remark` ,`legal` ,`path` ,`content` ,`enable` ,`create_time` ,`update_time` ,`is_del` from  %s where `is_del` != %d and `path` = ?", o.table, isDel)
	operations, err := o.fetch(ctx, query, path)
	if err != nil {
		log.Error("operationsRepository GetByPath fetch err : ", err.Error())
		return
	}
	if len(operations) > 0 {
		result = operations[0]
	}
	return
}

func (o *operationsRepository) GetByID(ctx context.Context, id int64) (result domain.FbOperations, err error) {
	query := fmt.Sprintf("SELECT `id` ,`method` ,`operation_type` ,`is_public` ,`remark` ,`legal` ,`path` ,`content` ,`enable` ,`create_time` ,`update_time` ,`is_del` from  %s where `is_del` != %d and `id` = ?", o.table, isDel)
	operations, err := o.fetch(ctx, query, id)
	if err != nil {
		log.Error("operationsRepository GetByID fetch err : ", err.Error())
		return
	}
	if len(operations) > 0 {
		result = operations[0]
	}
	return
}

func (o *operationsRepository) Delete(ctx context.Context, id int64) (result int64, err error) {
	query := fmt.Sprintf("DELETE from %s  where `id` = ? ", o.table)
	stmt, err := o.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error("operationsRepository Delete PrepareContext err : ", err.Error())
		return
	}
	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		log.Error("operationsRepository Delete ExecContext err : ", err.Error())
		return
	}
	result, err = res.RowsAffected()
	if err != nil {
		log.Error("operationsRepository Delete RowsAffected err : ", err.Error())
		return
	}

	return
}

func (o *operationsRepository) BatchDelete(ctx context.Context, ids string) (err error) {
	query := fmt.Sprintf("DELETE from %s  where `id` in (%s)) ", o.table, ids)
	stmt, err := o.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error("operationsRepository BatchDelete PrepareContext err : ", err.Error())
		return
	}
	_, err = stmt.ExecContext(ctx)
	if err != nil {
		log.Error("operationsRepository Delete ExecContext err : ", err.Error())
		return
	}

	return
}

func (o *operationsRepository) FindOperations(ctx context.Context) (result []domain.FbOperations, err error) {
	query := fmt.Sprintf("SELECT `id` ,`method` ,`operation_type` ,`is_public` ,`remark` ,`legal` ,`path` ,`content` ,`enable` ,`create_time` ,`update_time` ,`is_del` from  %s where `is_del` != %d ", o.table, isDel)
	result, err = o.fetch(ctx, query)
	if err != nil {
		log.Error("operationsRepository FindOperations fetch err : ", err.Error())
		return
	}
	return
}

func (o *operationsRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.FbOperations, err error) {
	rows, err := o.db.QueryContext(ctx, query, args...)
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

	result = make([]domain.FbOperations, 0)
	for rows.Next() {
		operation := domain.FbOperations{}
		err = rows.Scan(
			&operation.ID,
			&operation.Method,
			&operation.OperationType,
			&operation.IsPublic,
			&operation.Remark,
			&operation.Legal,
			&operation.Path,
			&operation.Content,
			&operation.Enable,
			&operation.CreateTime,
			&operation.UpdateTime,
			&operation.IsDel,
		)
		if err != nil {
			return
		}
		result = append(result, operation)
	}

	return result, nil
}

func GetUpdateSql(entity interface{}, tableName string) string {
	result := ""
	entityType := reflect.TypeOf(entity)
	entityValue := reflect.ValueOf(entity)

	// 反射判断该参数是否是结构体
	if entityValue.Kind() != reflect.Struct {
		return ""
	}
	result = fmt.Sprintf("UPDATE %s SET ", tableName)
	id := ""
	// 反射获取参数的每个字段,根据映射生成,所有字段模糊查询生成sql
	for i := 0; i < entityType.NumField(); i++ {
		if entityType.Field(i).Tag.Get("db") == "is_del" || entityType.Field(i).Tag.Get("db") == "update_time" ||
			entityType.Field(i).Tag.Get("db") == "create_time" {
			continue
		}
		if entityType.Field(i).Tag.Get("db") == "id" {
			id = utils.InterfaceToString(entityValue.Field(i).Interface())
			continue
		}

		// 获取值 如果是int64类型
		val := utils.InterfaceToString(entityValue.Field(i).Interface())
		if valInt64, err := strconv.Atoi(val); err == nil {
			result = fmt.Sprintf("%s %s = %v ,", result, entityType.Field(i).Tag.Get("db"), valInt64)
		}

		if !utils.Empty(entityValue.Field(i).Interface()) {
			result = fmt.Sprintf("%s %s = '%s' ,", result, entityType.Field(i).Tag.Get("db"), utils.InterfaceToString(entityValue.Field(i).Interface()))
		}
	}
	result = strings.TrimRight(result, ",")
	result = fmt.Sprintf("%s where id = %s", result, id)
	return result
}
