package sqllite

import (
	"context"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fire_boom/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

var mockStorageBucket = []domain.FbStorageBucket{
	{
		ID:         1,
		Name:       "1",
		Switch:     1,
		Config:     "1",
		CreateTime: dateFormatStr,
		UpdateTime: dateFormatStr,
		IsDel:      0,
	},
	{
		ID:         2,
		Name:       "2",
		Switch:     2,
		Config:     "2",
		CreateTime: dateFormatStr,
		UpdateTime: dateFormatStr,
		IsDel:      1,
	},
}

// TestDataSourceRepository_FindDataSources
func TestStorageBucketRepository_FindStorageBucket(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "switch", "config", "CreateTime", "UpdateTime", "isDel"}).
		AddRow(
			mockStorageBucket[0].ID,
			mockStorageBucket[0].Name,
			mockStorageBucket[0].Switch,
			mockStorageBucket[0].Config,
			mockStorageBucket[0].CreateTime,
			mockStorageBucket[0].UpdateTime,
			mockStorageBucket[0].IsDel,
		).
		AddRow(
			mockStorageBucket[1].ID,
			mockStorageBucket[1].Name,
			mockStorageBucket[1].Switch,
			mockStorageBucket[1].Config,
			mockStorageBucket[1].CreateTime,
			mockStorageBucket[1].UpdateTime,
			mockStorageBucket[1].IsDel,
		)
	s := NewStorageBucketRepository(db)
	query := fmt.Sprintf("SELECT `id`, `name`, `switch`, `config`, `create_time`, `update_time`, `is_del` from  %s where `is_del` != %d ", s.table, isDel)
	mock.ExpectQuery(query).WillReturnRows(rows)

	list, err := s.FindStorageBucket(context.TODO())
	assert.NotEmpty(t, list)
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestStorageBucketRepository_GetByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "name", "switch", "config", "CreateTime", "UpdateTime", "isDel"}).
		AddRow(
			mockStorageBucket[0].ID,
			mockStorageBucket[0].Name,
			mockStorageBucket[0].Switch,
			mockStorageBucket[0].Config,
			mockStorageBucket[0].CreateTime,
			mockStorageBucket[0].UpdateTime,
			mockStorageBucket[0].IsDel,
		)
	a := NewStorageBucketRepository(db)
	query := fmt.Sprintf("select `id`,`name`,`switch`,`config`,`create_time`,`update_time`,`is_del` from %s where `is_del` != 1 and `name` = ?", a.table)
	mock.ExpectQuery(query).WillReturnRows(rows)

	list, err := a.GetByName(context.TODO(), "1")
	assert.NotEmpty(t, list)
	assert.NoError(t, err)
}

func TestStorageBucketRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "name", "switch", "config", "CreateTime", "UpdateTime", "isDel"}).
		AddRow(
			mockStorageBucket[0].ID,
			mockStorageBucket[0].Name,
			mockStorageBucket[0].Switch,
			mockStorageBucket[0].Config,
			mockStorageBucket[0].CreateTime,
			mockStorageBucket[0].UpdateTime,
			mockStorageBucket[0].IsDel,
		)
	a := NewStorageBucketRepository(db)
	query := fmt.Sprintf("select `id`,`name`,`switch`,`config`,`create_time`,`update_time`,`is_del` from %s where `is_del` != 1 and `id` = ?", a.table)
	mock.ExpectQuery(query).WillReturnRows(rows)

	list, err := a.GetByID(context.TODO(), 1)
	assert.NotEmpty(t, list)
	assert.NoError(t, err)
}

func TestStorageBucketRepository_CheckExist(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "name", "switch", "config", "CreateTime", "UpdateTime", "isDel"}).
		AddRow(
			mockStorageBucket[0].ID,
			mockStorageBucket[0].Name,
			mockStorageBucket[0].Switch,
			mockStorageBucket[0].Config,
			mockStorageBucket[0].CreateTime,
			mockStorageBucket[0].UpdateTime,
			mockStorageBucket[0].IsDel,
		)
	a := NewStorageBucketRepository(db)
	query := fmt.Sprintf("select `id`,`name`,`switch`,`config`,`create_time`,`update_time`,`is_del` from %s where `is_del` != 1 and `id` != \\? and `name` = \\?", a.table)
	mock.ExpectQuery(query).WillReturnRows(rows)

	list, err := a.CheckExist(context.TODO(), &mockStorageBucket[0])
	assert.NotEmpty(t, list)
	assert.NoError(t, err)
}

func TestStorageBucketRepository_Update(t *testing.T) {
	t.Run("update one row", func(t *testing.T) {
		ds := mockStorageBucket[0]

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		a := NewStorageBucketRepository(db)
		query := fmt.Sprintf("UPDATE %s SET `name` = \\? , `switch` = \\? ,`config` = \\? WHERE `id` = \\?", a.table)
		prep := mock.ExpectPrepare(query)
		prep.ExpectExec().WithArgs(ds.Name, ds.Switch, ds.Config, ds.ID).WillReturnResult(sqlmock.NewResult(1, 1))

		_, err = a.Update(context.TODO(), &ds)
		assert.NoError(t, err)
	})

	t.Run("update two row", func(t *testing.T) {

		ds := mockStorageBucket
		ds[1].ID = 1

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		a := NewStorageBucketRepository(db)
		query := fmt.Sprintf("UPDATE %s SET `name` = \\? , `switch` = \\? ,`config` = \\? WHERE `id` = \\?", a.table)
		prep := mock.ExpectPrepare(query)
		prep.ExpectExec().WithArgs(ds[0].Name, ds[0].Switch, ds[0].Config, ds[0].ID).WillReturnResult(sqlmock.NewResult(2, 2))

		affect, err := a.Update(context.TODO(), &ds[0])
		if affect != 2 && err.Error() != "Weird  Behavior. Total Affected: 2" {
			t.Error("want affect=1 but affect=2")
		}
	})
}

func TestStorageBucketRepository_Delete(t *testing.T) {

	t.Run("delete one row", func(t *testing.T) {
		ds := mockStorageBucket[0]
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		a := NewStorageBucketRepository(db)
		query := fmt.Sprintf("UPDATE %s SET `is_del` = %d where `id` = \\? ", a.table, isDel)
		prep := mock.ExpectPrepare(query)
		prep.ExpectExec().WithArgs(ds.ID).WillReturnResult(sqlmock.NewResult(1, 1))

		_, err = a.Delete(context.TODO(), 1)
		assert.NoError(t, err)
	})

	t.Run("delete two row", func(t *testing.T) {
		ds := mockStorageBucket
		ds[1].ID = 1
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		a := NewStorageBucketRepository(db)
		query := fmt.Sprintf("UPDATE %s SET `is_del` = %d where `id` = \\? ", a.table, isDel)
		prep := mock.ExpectPrepare(query)
		prep.ExpectExec().WithArgs(1).WillReturnResult(sqlmock.NewResult(2, 2))

		affect, err := a.Delete(context.TODO(), 1)
		if affect != 2 && err.Error() != "Weird  Behavior. Total Affected: 2" {
			t.Error("want affect=1 but affect=2")
		}
	})
}

func TestStorageBucketRepository_Store(t *testing.T) {
	ds := &mockStorageBucket[0]
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	a := NewStorageBucketRepository(db)
	query := fmt.Sprintf("INSERT INTO %s \\(`name`, `switch`, `config`\\) values\\(\\?,\\?,\\?\\)", a.table)

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(ds.Name, ds.Switch, ds.Config).WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = a.Store(context.TODO(), ds)
	assert.NoError(t, err)
}
