package sqllite

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/fire_boom/domain"
	log "github.com/sirupsen/logrus"
)

type prismaRepository struct {
	db    *sql.DB
	table string
}

func NewPrismaRepository(Conn *sql.DB) *prismaRepository {
	return &prismaRepository{
		db:    Conn,
		table: "`fb_prisma`",
	}
}

func (p *prismaRepository) Create(ctx context.Context, prisma *domain.Prisma) (err error) {
	query := fmt.Sprintf("INSERT INTO %s (`name`, `file_id`) values(?,?)", p.table)
	stmt, err := p.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error("prismaRepository Create PrepareContext err : ", err.Error())
		return
	}
	_, err = stmt.ExecContext(ctx, prisma.Name, prisma.File.ID)
	if err != nil {
		log.Error("prismaRepository Create ExecContext err : ", err.Error())
		return
	}
	return
}

func (p *prismaRepository) Update(ctx context.Context, prisma *domain.Prisma) (err error) {
	query := fmt.Sprintf("UPDATE %s SET `name` = ? , `file_id` = ?  WHERE `id` = ?", p.table)
	stmt, err := p.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error("prismaRepository Update PrepareContext err : ", err.Error())
		return
	}
	_, err = stmt.ExecContext(ctx, prisma.Name, prisma.File.ID, prisma.ID)
	if err != nil {
		log.Error("prismaRepository Update ExecContext err : ", err.Error())
		return
	}
	return
}

func (p *prismaRepository) Fetch(ctx context.Context) (result []domain.Prisma, err error) {
	query := "select p.id prismaID,p.name prismaName,f.id fileID,f.name fileName,f.path filePath,p.create_time PrismaCreateTime,p.update_time PrismaUpdateTime,p.is_del PrismaIsDel from fb_prisma p left join fb_file f on p.file_id = f.id where p.is_del !=1 and f.is_del != 1"
	result, err = p.fetch(ctx, query)
	if err != nil {
		log.Error("prismaRepository Fetch fetch err : ", err.Error())
		return
	}
	return
}

func (p *prismaRepository) GetByID(ctx context.Context, id int64) (result domain.Prisma, err error) {
	query := "select p.id prismaID,p.name prismaName,f.id fileID,f.name fileName,f.path filePath,p.create_time PrismaCreateTime,p.update_time PrismaUpdateTime,p.is_del PrismaIsDel from fb_prisma p left join fb_file f on p.file_id = f.id where p.is_del !=1 and f.is_del != 1 and f.id = ?"
	var queryOne []domain.Prisma
	queryOne, err = p.fetch(ctx, query, id)
	if err != nil {
		log.Error("prismaRepository GetByID fetch err : ", err.Error())
		return
	}
	if len(queryOne) > 0 {
		result = queryOne[0]
	}
	return
}

func (p *prismaRepository) CheckExist(ctx context.Context, prisma *domain.Prisma) (result domain.Prisma, err error) {
	query := "select p.id prismaID,p.name prismaName,f.id fileID,f.name fileName,f.path filePath,p.create_time PrismaCreateTime,p.update_time PrismaUpdateTime,p.is_del PrismaIsDel from fb_prisma p left join fb_file f on p.file_id = f.id where p.is_del !=1 and f.is_del != 1 and f.id != ? and f.name = ?"
	var queryOne []domain.Prisma
	queryOne, err = p.fetch(ctx, query, prisma.ID, prisma.Name)
	if err != nil {
		log.Error("prismaRepository CheckExist fetch err : ", err.Error())
		return
	}
	if len(queryOne) > 0 {
		result = queryOne[0]
	}
	return
}

func (p *prismaRepository) GetByName(ctx context.Context, name string) (result domain.Prisma, err error) {
	query := "select p.id prismaID,p.name prismaName,f.id fileID,f.name fileName,f.path filePath,p.create_time PrismaCreateTime,p.update_time PrismaUpdateTime,p.is_del PrismaIsDel from fb_prisma p left join fb_file f on p.file_id = f.id where p.is_del !=1 and f.is_del != 1 and f.name = ?"
	var queryOne []domain.Prisma
	queryOne, err = p.fetch(ctx, query, name)
	if err != nil {
		log.Error("prismaRepository GetByName fetch err : ", err.Error())
		return
	}
	if len(queryOne) > 0 {
		result = queryOne[0]
	}
	return
}

func (p *prismaRepository) Delete(ctx context.Context, id int64) (err error) {
	query := fmt.Sprintf("UPDATE %s SET `is_del` = %d where `id` = ? ", p.table, isDel)
	stmt, err := p.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error("prismaRepository Delete PrepareContext err : ", err.Error())
		return
	}
	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		log.Error("prismaRepository Delete ExecContext err : ", err.Error())
		return
	}

	return
}

func (p *prismaRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Prisma, err error) {
	rows, err := p.db.QueryContext(ctx, query)

	if err != nil {
		log.Error("prismaRepository fetch QueryContext err : ", err.Error())
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			err = errRow
			return
		}
	}()

	result = make([]domain.Prisma, 0)
	for rows.Next() {
		prisma := domain.Prisma{}
		file := domain.File{}
		err = rows.Scan(
			&prisma.ID,
			&prisma.Name,
			&file.ID,
			&file.Name,
			&file.Path,
			&prisma.CreateTime,
			&prisma.UpdateTime,
			&prisma.IsDel,
		)
		prisma.File = file
		if err != nil {
			return
		}
		result = append(result, prisma)
	}

	return result, nil
}
