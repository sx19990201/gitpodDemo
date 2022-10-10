package sqllite

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fire_boom/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	dateFormatStr = "2006-01-02 15:04:05"
)

var mockDataSource = []domain.FbDataSource{
	{
		ID:         1,
		Name:       "1",
		SourceType: 1,
		Config:     "1",
		CreateTime: sql.NullString{String: dateFormatStr},
		UpdateTime: sql.NullString{String: dateFormatStr},
		IsDel:      0,
	},
	{
		ID:         2,
		Name:       "2",
		SourceType: 2,
		Config:     "2",
		CreateTime: sql.NullString{String: "2022-06-23 00:00:00"},
		UpdateTime: sql.NullString{String: "2022-06-23 00:00:00"},
		IsDel:      1,
	},
}

// TestDataSourceRepository_FindDataSources
func TestDataSourceRepository_FindDataSources(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "sourceType", "config", "switch", "create_time", "update_time", "isDel"}).
		AddRow(
			mockDataSource[0].ID,
			mockDataSource[0].Name,
			mockDataSource[0].SourceType,
			mockDataSource[0].Config,
			mockDataSource[0].Switch,
			mockDataSource[0].CreateTime,
			mockDataSource[0].UpdateTime,
			mockDataSource[0].IsDel,
		).
		AddRow(
			mockDataSource[1].ID,
			mockDataSource[1].Name,
			mockDataSource[1].SourceType,
			mockDataSource[1].Config,
			mockDataSource[1].Switch,
			mockDataSource[1].CreateTime,
			mockDataSource[1].UpdateTime,
			mockDataSource[1].IsDel,
		)

	query := "SELECT `id`, `name`, `source_type`, `config`,`switch`, `create_time`, `update_time`, `is_del` from `fb_data_source` where `is_del` != 1 "
	mock.ExpectQuery(query).WillReturnRows(rows)

	a := NewDataSourceRepository(db)
	list, err := a.FindDataSources(context.TODO())
	assert.NotEmpty(t, list)
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestDataSourceRepository_GetByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "name", "sourceType", "config", "switch", "create_time", "update_time", "isDel"}).
		AddRow(
			mockDataSource[0].ID,
			mockDataSource[0].Name,
			mockDataSource[0].SourceType,
			mockDataSource[0].Config,
			mockDataSource[0].Switch,
			mockDataSource[0].CreateTime,
			mockDataSource[0].UpdateTime,
			mockDataSource[0].IsDel,
		)
	a := NewDataSourceRepository(db)
	query := fmt.Sprintf("select `id`,`name`,`source_type`,`config`,`switch`,`create_time`,`update_time`,`is_del` from %s where `is_del` != 1 and `name` = ?", a.table)
	mock.ExpectQuery(query).WillReturnRows(rows)

	list, err := a.GetByName(context.TODO(), "1")
	assert.NotEmpty(t, list)
	assert.NoError(t, err)
}

func TestDataSourceRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "name", "sourceType", "config", "switch", "create_time", "update_time", "isDel"}).
		AddRow(
			mockDataSource[0].ID,
			mockDataSource[0].Name,
			mockDataSource[0].SourceType,
			mockDataSource[0].Config,
			mockDataSource[0].Switch,
			mockDataSource[0].CreateTime,
			mockDataSource[0].UpdateTime,
			mockDataSource[0].IsDel,
		)
	a := NewDataSourceRepository(db)
	query := fmt.Sprintf("select `id`,`name`,`source_type`,`config`,`switch`,`create_time`,`update_time`,`is_del` from %s where `is_del` != 1 and `id` = ?", a.table)
	mock.ExpectQuery(query).WillReturnRows(rows)

	list, err := a.GetByID(context.TODO(), 1)
	assert.NotEmpty(t, list)
	assert.NoError(t, err)
}

func TestDataSourceRepository_CheckExist(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "name", "sourceType", "config", "switch", "create_time", "update_time", "isDel"}).
		AddRow(
			mockDataSource[0].ID,
			mockDataSource[0].Name,
			mockDataSource[0].SourceType,
			mockDataSource[0].Config,
			mockDataSource[0].Switch,
			mockDataSource[0].CreateTime,
			mockDataSource[0].UpdateTime,
			mockDataSource[0].IsDel,
		)
	a := NewDataSourceRepository(db)
	query := fmt.Sprintf("select `id`,`name`,`source_type`,`config`,`switch`,`create_time`,`update_time`,`is_del` from %s where `is_del` != 1 and `id` != \\? and `name` = \\?", a.table)
	mock.ExpectQuery(query).WillReturnRows(rows)

	list, err := a.CheckExist(context.TODO(), &mockDataSource[0])
	assert.NotEmpty(t, list)
	assert.NoError(t, err)
}

func TestDataSourceRepository_Update(t *testing.T) {
	t.Run("update one row", func(t *testing.T) {
		ds := mockDataSource[0]

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		a := NewDataSourceRepository(db)
		query := fmt.Sprintf("UPDATE %s SET `name` = \\? , `source_type` = \\? ,`config` = \\?,`switch`= \\? WHERE `id` = \\?", a.table)
		prep := mock.ExpectPrepare(query)
		prep.ExpectExec().WithArgs(ds.Name, ds.SourceType, ds.Config, ds.Switch, ds.ID).WillReturnResult(sqlmock.NewResult(1, 1))

		_, err = a.Update(context.TODO(), &ds)
		assert.NoError(t, err)
	})

	t.Run("update two row", func(t *testing.T) {

		ds := mockDataSource
		ds[1].ID = 1

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		a := NewDataSourceRepository(db)
		query := fmt.Sprintf("UPDATE %s SET `name` = \\? , `source_type` = \\? ,`config` = \\?,`switch`= \\? WHERE `id` = \\?", a.table)
		prep := mock.ExpectPrepare(query)
		prep.ExpectExec().WithArgs(ds[0].Name, ds[0].SourceType, ds[0].Config, ds[0].Switch, ds[0].ID).WillReturnResult(sqlmock.NewResult(2, 2))

		affect, err := a.Update(context.TODO(), &ds[0])
		if affect != 2 && err.Error() != "Weird  Behavior. Total Affected: 2" {
			t.Error("want affect=1 but affect=2")
		}
	})
}

func TestDataSourceRepository_Delete(t *testing.T) {

	t.Run("delete one row", func(t *testing.T) {
		ds := mockDataSource[0]
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		a := NewDataSourceRepository(db)
		query := fmt.Sprintf("UPDATE %s SET `is_del` = %d where `id` = \\? ", a.table, isDel)
		prep := mock.ExpectPrepare(query)
		prep.ExpectExec().WithArgs(ds.ID).WillReturnResult(sqlmock.NewResult(1, 1))

		_, err = a.Delete(context.TODO(), 1)
		assert.NoError(t, err)
	})

	t.Run("delete two row", func(t *testing.T) {
		ds := mockDataSource
		ds[1].ID = 1
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		a := NewDataSourceRepository(db)
		query := fmt.Sprintf("UPDATE %s SET `is_del` = %d where `id` = \\? ", a.table, isDel)
		prep := mock.ExpectPrepare(query)
		prep.ExpectExec().WithArgs(1).WillReturnResult(sqlmock.NewResult(2, 2))

		affect, err := a.Delete(context.TODO(), 1)
		if affect != 2 && err.Error() != "Weird  Behavior. Total Affected: 2" {
			t.Error("want affect=1 but affect=2")
		}
	})

}

func TestDataSourceRepository_Store(t *testing.T) {
	ds := &mockDataSource[0]
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	a := NewDataSourceRepository(db)
	query := fmt.Sprintf("INSERT INTO %s \\(`name`, `source_type`, `config`,`switch`\\) values\\(\\?,\\?,\\?,\\?\\)", a.table)

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(ds.Name, ds.SourceType, ds.Config, ds.Switch).WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = a.Store(context.TODO(), ds)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), ds.ID)
}
