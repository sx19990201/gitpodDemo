package sqllite

import (
	"context"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fire_boom/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

var mockPrisma = []domain.Prisma{
	{
		ID:   1,
		Name: "1",
		File: domain.File{
			ID:   "1",
			Name: "1",
			Path: "1",
		},
		IsDel: 0,
	}, {
		ID:   2,
		Name: "2",
		File: domain.File{
			ID:   "2",
			Name: "2",
			Path: "2",
		},
		IsDel: 0,
	},
}

func TestPrismaRepository_Create(t *testing.T) {
	mockP := &mockPrisma[0]
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	roleRepo := NewPrismaRepository(db)
	query := fmt.Sprintf("INSERT INTO %s \\(`name`, `file_id`\\) values\\(\\?,\\?\\)", roleRepo.table)

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(mockP.Name, mockP.File.ID).WillReturnResult(sqlmock.NewResult(1, 1))

	err = roleRepo.Create(context.TODO(), mockP)
	assert.NoError(t, err)
}

func TestPrismaRepository_Update(t *testing.T) {
	t.Run("update one row", func(t *testing.T) {
		mockP := mockPrisma[0]

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		a := NewPrismaRepository(db)
		query := fmt.Sprintf("UPDATE %s SET `name` = \\? , `file_id` = \\? WHERE `id` = \\?", a.table)
		prep := mock.ExpectPrepare(query)
		prep.ExpectExec().WithArgs(mockP.Name, mockP.File.ID, mockP.ID).WillReturnResult(sqlmock.NewResult(1, 1))

		err = a.Update(context.TODO(), &mockP)
		assert.NoError(t, err)
	})

	t.Run("update two row", func(t *testing.T) {

		mockP := mockPrisma
		mockP[1].ID = 1

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		a := NewPrismaRepository(db)
		query := fmt.Sprintf("UPDATE %s SET `name` = \\? , `file_id` = \\? WHERE `id` = \\?", a.table)
		prep := mock.ExpectPrepare(query)
		prep.ExpectExec().WithArgs(mockP[0].Name, mockP[0].File.ID, mockP[0].ID).WillReturnResult(sqlmock.NewResult(2, 2))

		err = a.Update(context.TODO(), &mockP[0])
		assert.NoError(t, err)
	})
}

func TestPrismaRepository_Delete(t *testing.T) {
	t.Run("delete one row", func(t *testing.T) {
		ds := mockPrisma[0]
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		a := NewPrismaRepository(db)
		query := fmt.Sprintf("UPDATE %s SET `is_del` = %d where `id` = \\? ", a.table, isDel)
		prep := mock.ExpectPrepare(query)
		prep.ExpectExec().WithArgs(ds.ID).WillReturnResult(sqlmock.NewResult(1, 1))

		err = a.Delete(context.TODO(), 1)
		assert.NoError(t, err)
	})

	t.Run("delete two row", func(t *testing.T) {
		ds := mockPrisma
		ds[1].ID = 1
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		a := NewPrismaRepository(db)
		query := fmt.Sprintf("UPDATE %s SET `is_del` = %d where `id` = \\? ", a.table, isDel)
		prep := mock.ExpectPrepare(query)
		prep.ExpectExec().WithArgs(1).WillReturnResult(sqlmock.NewResult(2, 2))

		err = a.Delete(context.TODO(), 1)
		assert.NoError(t, err)
	})
}

func TestPrismaRepository_Fetch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"prismaID", "prismaName", "fileID", "fileName", "filePath", "PrismaCreateTime", "PrismaUpdateTime", "PrismaIsDel"}).
		AddRow(
			mockPrisma[0].ID,
			mockPrisma[0].Name,
			mockPrisma[0].File.ID,
			mockPrisma[0].File.Name,
			mockPrisma[0].File.Path,
			mockPrisma[0].CreateTime,
			mockPrisma[0].UpdateTime,
			mockPrisma[0].IsDel,
		).
		AddRow(
			mockPrisma[1].ID,
			mockPrisma[1].Name,
			mockPrisma[1].File.ID,
			mockPrisma[1].File.Name,
			mockPrisma[1].File.Path,
			mockPrisma[1].CreateTime,
			mockPrisma[1].UpdateTime,
			mockPrisma[1].IsDel,
		)
	a := NewPrismaRepository(db)

	query := "select p.id prismaID,p.name prismaName,f.id fileID,f.name fileName,f.path filePath,p.create_time PrismaCreateTime,p.update_time PrismaUpdateTime,p.is_del PrismaIsDel from fb_prisma p left join fb_file f on p.file_id = f.id where p.is_del !=1 and f.is_del != 1"
	mock.ExpectQuery(query).WillReturnRows(rows)
	list, err := a.Fetch(context.TODO())
	assert.NotEmpty(t, list)
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestPrismaRepository_GetByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"prismaID", "prismaName", "fileID", "fileName", "filePath", "PrismaCreateTime", "PrismaUpdateTime", "PrismaIsDel"}).
		AddRow(
			mockPrisma[0].ID,
			mockPrisma[0].Name,
			mockPrisma[0].File.ID,
			mockPrisma[0].File.Name,
			mockPrisma[0].File.Path,
			mockPrisma[0].CreateTime,
			mockPrisma[0].UpdateTime,
			mockPrisma[0].IsDel,
		)

	a := NewPrismaRepository(db)
	query := "select p.id prismaID,p.name prismaName,f.id fileID,f.name fileName,f.path filePath,p.create_time PrismaCreateTime,p.update_time PrismaUpdateTime,p.is_del PrismaIsDel from fb_prisma p left join fb_file f on p.file_id = f.id where p.is_del !=1 and f.is_del != 1 and f.name = ?"
	mock.ExpectQuery(query).WillReturnRows(rows)

	list, err := a.GetByName(context.TODO(), "1")
	assert.NotEmpty(t, list)
	assert.NoError(t, err)
}

func TestPrismaRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"prismaID", "prismaName", "fileID", "fileName", "filePath", "PrismaCreateTime", "PrismaUpdateTime", "PrismaIsDel"}).
		AddRow(
			mockPrisma[0].ID,
			mockPrisma[0].Name,
			mockPrisma[0].File.ID,
			mockPrisma[0].File.Name,
			mockPrisma[0].File.Path,
			mockPrisma[0].CreateTime,
			mockPrisma[0].UpdateTime,
			mockPrisma[0].IsDel,
		)

	a := NewPrismaRepository(db)
	query := "select p.id prismaID,p.name prismaName,f.id fileID,f.name fileName,f.path filePath,p.create_time PrismaCreateTime,p.update_time PrismaUpdateTime,p.is_del PrismaIsDel from fb_prisma p left join fb_file f on p.file_id = f.id where p.is_del !=1 and f.is_del != 1 and f.id = ?"
	mock.ExpectQuery(query).WillReturnRows(rows)

	list, err := a.GetByID(context.TODO(), 1)
	assert.NotEmpty(t, list)
	assert.NoError(t, err)
}

func TestPrismaRepository_CheckExist(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"prismaID", "prismaName", "fileID", "fileName", "filePath", "PrismaCreateTime", "PrismaUpdateTime", "PrismaIsDel"}).
		AddRow(
			mockPrisma[0].ID,
			mockPrisma[0].Name,
			mockPrisma[0].File.ID,
			mockPrisma[0].File.Name,
			mockPrisma[0].File.Path,
			mockPrisma[0].CreateTime,
			mockPrisma[0].UpdateTime,
			mockPrisma[0].IsDel,
		)
	a := NewPrismaRepository(db)
	query := "select p.id prismaID,p.name prismaName,f.id fileID,f.name fileName,f.path filePath,p.create_time PrismaCreateTime,p.update_time PrismaUpdateTime,p.is_del PrismaIsDel from fb_prisma p left join fb_file f on p.file_id = f.id where p.is_del !=1 and f.is_del != 1 and f.id != \\? and f.name = \\?"
	mock.ExpectQuery(query).WillReturnRows(rows)

	list, err := a.CheckExist(context.TODO(), &mockPrisma[0])
	assert.NotEmpty(t, list)
	assert.NoError(t, err)
}
