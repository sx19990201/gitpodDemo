package http

import (
	"database/sql"
	"fmt"
	"github.com/fire_boom/domain"
	"github.com/fire_boom/pkg/repository/sqllite"
	"github.com/fire_boom/pkg/usecase"
	"github.com/fire_boom/utils"
	"github.com/labstack/echo/v4"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

type fileHandler struct {
	FileUseCase domain.FileUseCase
}

func NewFileHandler(e *echo.Echo, us domain.FileUseCase) {
	handler := &fileHandler{FileUseCase: us}
	v1 := e.Group("/api/v1")
	{
		dataSource := v1.Group("/file")
		{
			dataSource.POST("/upload", handler.Store)
			dataSource.DELETE("", handler.Delete)
		}
	}
}

func InitFileHandler(e *echo.Echo, db *sql.DB) {
	fileRepo := sqllite.NewFileRepository(db)
	timeoutContext := time.Duration(utils.GetTimeOut()) * time.Second
	au := usecase.NewFileUseCase(fileRepo, timeoutContext)
	NewFileHandler(e, au)
}

// Store @Title Store
// @Description 添加保存文件
// @Accept  json
// @Tags  file
// @Param file formData string true "file"
// @Success 200 "添加成功"
// @Failure 400	"添加失败"
// @Router /api/v1/file/upload  [POST]
func (f *fileHandler) Store(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.FileCode, domain.ParamErr))
	}
	fileUUID := utils.GenerateUUID()
	filePath := fmt.Sprintf("%s/%s", utils.GetOASFilePath(), fileUUID)
	err = uploadWriteFile(file, filePath)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.FileCode, domain.FileWriteErr))
	}

	ctx := c.Request().Context()
	_, err = f.FileUseCase.Store(ctx, &domain.File{
		ID:   fileUUID,
		Name: file.Filename,
		Path: filePath,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.FileCode, err))
	}

	return c.JSON(http.StatusOK, SuccessWriteResult(fileUUID))
}

// Delete @Title Delete
// @Description 删除文件
// @Accept  json
// @Tags  file
// @Param id formData string true "id"
// @Success 200 "删除成功"
// @Failure 400	"删除失败"
// @Router /api/v1/file  [Delete]
func (f *fileHandler) Delete(c echo.Context) error {
	param := struct {
		UUID string `json:"id"`
		Path string `json:"path"`
	}{}
	err := c.Bind(&param)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.FileCode, domain.ParamErr))
	}

	filePath := fmt.Sprintf("%s/%s", param.Path, param.UUID)
	err = os.Remove(filePath)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.FileCode, domain.FileDeleteErr))
	}
	ctx := c.Request().Context()
	err = f.FileUseCase.Delete(ctx, param.UUID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.FileCode, domain.FileDeleteErr))
	}

	return c.JSON(http.StatusOK, SuccessResult())
}

func uploadWriteFile(file *multipart.FileHeader, path string) error {
	src, err := file.Open()
	if err != nil {
		return domain.FileOpenErr
	}
	defer src.Close()

	dst, err := os.Create(path)
	if err != nil {
		return domain.FileCreateErr
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return domain.FileCopyErr
	}
	return nil
}
