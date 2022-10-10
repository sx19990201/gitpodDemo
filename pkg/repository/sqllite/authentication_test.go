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

var mockAuthentication = []domain.FbAuthentication{
	{
		ID:           1,
		Name:         "1",
		AuthSupplier: "1",
		Config:       "1",
		CreateTime:   sql.NullString{String: dateFormatStr},
		UpdateTime:   sql.NullString{String: dateFormatStr},
		IsDel:        0,
	}, {
		ID:           2,
		Name:         "2",
		AuthSupplier: "2",
		Config:       "2",
		CreateTime:   sql.NullString{String: dateFormatStr},
		UpdateTime:   sql.NullString{String: dateFormatStr},
		IsDel:        0,
	},
}

func TestAuthenticationRepository_Store(t *testing.T) {
	ds := &mockAuthentication[0]
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	authRepo := NewAuthenticationRepository(db)
	query := fmt.Sprintf("INSERT INTO %s \\(`name`, `auth_supplier`, `config`, `switch_state`\\) values\\(\\?,\\?,\\?,\\?\\)", authRepo.table)

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(ds.Name, ds.AuthSupplier, ds.Config, ds.SwitchState).WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = authRepo.Store(context.TODO(), ds)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), ds.ID)
}

func TestAuthenticationRepository_Update(t *testing.T) {
	t.Run("update one row", func(t *testing.T) {
		mockAuth := mockAuthentication[0]

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		a := NewAuthenticationRepository(db)
		query := fmt.Sprintf("UPDATE %s SET `name` = \\? , `auth_supplier` = \\? ,`config` = \\?, `switch_state` = \\? WHERE `id` = \\?", a.table)
		prep := mock.ExpectPrepare(query)
		prep.ExpectExec().WithArgs(mockAuth.Name, mockAuth.AuthSupplier, mockAuth.Config, mockAuth.SwitchState, mockAuth.ID).WillReturnResult(sqlmock.NewResult(1, 1))

		_, err = a.Update(context.TODO(), &mockAuth)
		assert.NoError(t, err)
	})

	t.Run("update two row", func(t *testing.T) {

		mockAuth := mockAuthentication
		mockAuth[1].ID = 1

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		a := NewAuthenticationRepository(db)
		query := fmt.Sprintf("UPDATE %s SET `name` = \\? , `auth_supplier` = \\? ,`config` = \\?, `switch_state` = \\? WHERE `id` = \\?", a.table)
		prep := mock.ExpectPrepare(query)
		prep.ExpectExec().WithArgs(mockAuth[0].Name, mockAuth[0].AuthSupplier, mockAuth[0].Config, mockAuth[0].SwitchState, mockAuth[0].ID).WillReturnResult(sqlmock.NewResult(2, 2))

		affect, err := a.Update(context.TODO(), &mockAuth[0])
		if affect != 2 && err.Error() != "Weird  Behavior. Total Affected: 2" {
			t.Error("want affect=1 but affect=2")
		}
	})
}

func TestAuthenticationRepository_GetByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "name", "authSupplier", "switchState", "config", "create_time", "update_time", "isDel"}).
		AddRow(
			mockAuthentication[0].ID,
			mockAuthentication[0].Name,
			mockAuthentication[0].AuthSupplier,
			mockAuthentication[0].SwitchState,
			mockAuthentication[0].Config,
			mockAuthentication[0].CreateTime,
			mockAuthentication[0].UpdateTime,
			mockAuthentication[0].IsDel,
		)
	a := NewAuthenticationRepository(db)
	query := fmt.Sprintf("select `id`,`name`,`auth_supplier`,`switch_state`,`config`,`create_time`,`update_time`,`is_del` from %s where `is_del` != 1 and `name` = ?", a.table)
	mock.ExpectQuery(query).WillReturnRows(rows)

	list, err := a.GetByName(context.TODO(), "1")
	assert.NotEmpty(t, list)
	assert.NoError(t, err)
}

func TestAuthenticationRepository_CheckExist(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "name", "authSupplier", "switchState", "config", "create_time", "update_time", "isDel"}).
		AddRow(
			mockAuthentication[0].ID,
			mockAuthentication[0].Name,
			mockAuthentication[0].AuthSupplier,
			mockAuthentication[0].SwitchState,
			mockAuthentication[0].Config,
			mockAuthentication[0].CreateTime,
			mockAuthentication[0].UpdateTime,
			mockAuthentication[0].IsDel,
		)
	a := NewAuthenticationRepository(db)
	query := fmt.Sprintf("select `id`,`name`,`auth_supplier`,`switch_state`,`config`,`create_time`,`update_time`,`is_del` from %s where `is_del` != 1 and `id` != \\? and `name` = \\?", a.table)
	mock.ExpectQuery(query).WillReturnRows(rows)

	list, err := a.CheckExist(context.TODO(), &mockAuthentication[0])
	assert.NotEmpty(t, list)
	assert.NoError(t, err)
}

func TestAuthenticationRepository_Delete(t *testing.T) {
	t.Run("delete one row", func(t *testing.T) {
		ds := mockAuthentication[0]
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		a := NewAuthenticationRepository(db)
		query := fmt.Sprintf("UPDATE %s SET `is_del` = %d where `id` = \\? ", a.table, isDel)
		prep := mock.ExpectPrepare(query)
		prep.ExpectExec().WithArgs(ds.ID).WillReturnResult(sqlmock.NewResult(1, 1))

		_, err = a.Delete(context.TODO(), 1)
		assert.NoError(t, err)
	})

	t.Run("delete two row", func(t *testing.T) {
		ds := mockAuthentication
		ds[1].ID = 1
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		a := NewAuthenticationRepository(db)
		query := fmt.Sprintf("UPDATE %s SET `is_del` = %d where `id` = \\? ", a.table, isDel)
		prep := mock.ExpectPrepare(query)
		prep.ExpectExec().WithArgs(1).WillReturnResult(sqlmock.NewResult(2, 2))

		affect, err := a.Delete(context.TODO(), 1)
		if affect != 2 && err.Error() != "Weird  Behavior. Total Affected: 2" {
			t.Error("want affect=1 but affect=2")
		}
	})
}

func TestAuthenticationRepository_FindAuthentication(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "authSupplier", "switchState", "config", "create_time", "update_time", "isDel"}).
		AddRow(
			mockAuthentication[0].ID,
			mockAuthentication[0].Name,
			mockAuthentication[0].AuthSupplier,
			mockAuthentication[0].SwitchState,
			mockAuthentication[0].Config,
			mockAuthentication[0].CreateTime,
			mockAuthentication[0].UpdateTime,
			mockAuthentication[0].IsDel,
		).
		AddRow(
			mockAuthentication[1].ID,
			mockAuthentication[1].Name,
			mockAuthentication[1].AuthSupplier,
			mockAuthentication[1].SwitchState,
			mockAuthentication[1].Config,
			mockAuthentication[1].CreateTime,
			mockAuthentication[1].UpdateTime,
			mockAuthentication[1].IsDel,
		)
	a := NewAuthenticationRepository(db)

	query := fmt.Sprintf("SELECT `id`, `name`, `auth_supplier`,`switch_state`, `config`, `create_time`, `update_time`, `is_del` from %s where `is_del` != %d ", a.table, isDel)
	mock.ExpectQuery(query).WillReturnRows(rows)
	list, err := a.FindAuthentication(context.TODO())
	assert.NotEmpty(t, list)
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}
