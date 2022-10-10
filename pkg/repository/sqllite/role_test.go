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

var mockRoles = []domain.FbRole{
	{
		ID:         1,
		Code:       "1",
		Remark:     "1",
		CreateTime: sql.NullString{String: dateFormatStr},
		UpdateTime: sql.NullString{String: dateFormatStr},
		IsDel:      0,
	}, {
		ID:         2,
		Code:       "2",
		Remark:     "2",
		CreateTime: sql.NullString{String: dateFormatStr},
		UpdateTime: sql.NullString{String: dateFormatStr},
		IsDel:      0,
	},
}

func TestRoleRepository_Store(t *testing.T) {
	mockR := &mockRoles[0]
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	roleRepo := NewRoleRepository(db)
	query := fmt.Sprintf("INSERT INTO %s \\(`code`, `remark`\\) values\\(\\?,\\?\\)", roleRepo.table)

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(mockR.Code, mockR.Remark).WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = roleRepo.Store(context.TODO(), mockR)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), mockR.ID)
}

func TestRoleRepository_Update(t *testing.T) {
	t.Run("update one row", func(t *testing.T) {
		mockR := mockRoles[0]

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		a := NewRoleRepository(db)
		query := fmt.Sprintf("UPDATE %s SET `code` = \\? , `remark` = \\? WHERE `id` = \\?", a.table)
		prep := mock.ExpectPrepare(query)
		prep.ExpectExec().WithArgs(mockR.Code, mockR.Remark, mockR.ID).WillReturnResult(sqlmock.NewResult(1, 1))

		_, err = a.Update(context.TODO(), &mockR)
		assert.NoError(t, err)
	})

	t.Run("update two row", func(t *testing.T) {

		mockR := mockRoles
		mockR[1].ID = 1

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		a := NewRoleRepository(db)
		query := fmt.Sprintf("UPDATE %s SET `code` = \\? , `remark` = \\? WHERE `id` = \\?", a.table)
		prep := mock.ExpectPrepare(query)
		prep.ExpectExec().WithArgs(mockR[0].Code, mockR[0].Remark, mockR[0].ID).WillReturnResult(sqlmock.NewResult(2, 2))

		affect, err := a.Update(context.TODO(), &mockR[0])
		if affect != 2 && err.Error() != "Weird  Behavior. Total Affected: 2" {
			t.Error("want affect=1 but affect=2")
		}
	})
}

func TestRoleRepository_GetByCode(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "code", "remark", "CreateTime", "UpdateTime", "isDel"}).
		AddRow(
			mockRoles[0].ID,
			mockRoles[0].Code,
			mockRoles[0].Remark,
			mockRoles[0].CreateTime,
			mockRoles[0].UpdateTime,
			mockRoles[0].IsDel,
		)
	a := NewRoleRepository(db)
	query := fmt.Sprintf("select `id`,`code`,`remark`,`create_time`,`update_time`,`is_del` from %s where `is_del` != 1 and `code` = ?", a.table)
	mock.ExpectQuery(query).WillReturnRows(rows)

	list, err := a.GetByCode(context.TODO(), "1")
	assert.NotEmpty(t, list)
	assert.NoError(t, err)
}

func TestRoleRepository_CheckExist(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "code", "remark", "CreateTime", "UpdateTime", "isDel"}).
		AddRow(
			mockRoles[0].ID,
			mockRoles[0].Code,
			mockRoles[0].Remark,
			mockRoles[0].CreateTime,
			mockRoles[0].UpdateTime,
			mockRoles[0].IsDel,
		)
	a := NewRoleRepository(db)
	query := fmt.Sprintf("select `id`,`code`,`remark`,`create_time`,`update_time`,`is_del` from %s where `is_del` != 1 and id != \\? and `code` = \\?", a.table)
	mock.ExpectQuery(query).WillReturnRows(rows)

	list, err := a.CheckExist(context.TODO(), &mockRoles[0])
	assert.NotEmpty(t, list)
	assert.NoError(t, err)
}

func TestRoleRepository_Delete(t *testing.T) {
	t.Run("delete one row", func(t *testing.T) {
		ds := mockRoles[0]
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		a := NewRoleRepository(db)
		query := fmt.Sprintf("UPDATE %s SET `is_del` = %d where `id` = \\? ", a.table, isDel)
		prep := mock.ExpectPrepare(query)
		prep.ExpectExec().WithArgs(ds.ID).WillReturnResult(sqlmock.NewResult(1, 1))

		_, err = a.Delete(context.TODO(), 1)
		assert.NoError(t, err)
	})

	t.Run("delete two row", func(t *testing.T) {
		ds := mockRoles
		ds[1].ID = 1
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		a := NewRoleRepository(db)
		query := fmt.Sprintf("UPDATE %s SET `is_del` = %d where `id` = \\? ", a.table, isDel)
		prep := mock.ExpectPrepare(query)
		prep.ExpectExec().WithArgs(1).WillReturnResult(sqlmock.NewResult(2, 2))

		affect, err := a.Delete(context.TODO(), 1)
		if affect != 2 && err.Error() != "Weird  Behavior. Total Affected: 2" {
			t.Error("want affect=1 but affect=2")
		}
	})
}

func TestRoleRepository_FindRoles(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "code", "remark", "CreateTime", "UpdateTime", "isDel"}).
		AddRow(
			mockRoles[0].ID,
			mockRoles[1].Code,
			mockRoles[1].Remark,
			mockRoles[0].CreateTime,
			mockRoles[0].UpdateTime,
			mockRoles[0].IsDel,
		).
		AddRow(
			mockRoles[1].ID,
			mockRoles[1].Code,
			mockRoles[1].Remark,
			mockRoles[1].CreateTime,
			mockRoles[1].UpdateTime,
			mockRoles[1].IsDel,
		)
	a := NewRoleRepository(db)

	query := fmt.Sprintf("SELECT `id`, `code`, `remark`, `create_time`, `update_time`, `is_del` from  %s where `is_del` != %d ", a.table, isDel)
	mock.ExpectQuery(query).WillReturnRows(rows)
	list, err := a.FindRoles(context.TODO())
	assert.NotEmpty(t, list)
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}
