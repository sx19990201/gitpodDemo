package http

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/fire_boom/domain"
	"github.com/fire_boom/pkg/repository/sqllite"
	"github.com/fire_boom/pkg/usecase"
	"github.com/fire_boom/utils"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type S3UploadClientHandler struct {
	S3UploadClientUseCase domain.S3UploadClientUseCase
	StorageBucketUseCase  domain.StorageBucketUseCase
}

func InitS3UploadClientUseCaseRouter(e *echo.Echo, db *sql.DB) {
	s3Repo := sqllite.NewS3UploadClientRepository()
	storageRepo := sqllite.NewStorageBucketRepository(db)
	timeoutContext := time.Duration(utils.GetTimeOut()) * time.Second
	s3UC := usecase.NewS3UploadClientUseCase(s3Repo, timeoutContext)
	storageUC := usecase.NewStorageBucketUseCase(storageRepo, timeoutContext)
	NewS3UploadClientHandler(e, s3UC, storageUC)
}

func NewS3UploadClientHandler(e *echo.Echo, sUseCase domain.S3UploadClientUseCase, storageUC domain.StorageBucketUseCase) {
	handler := &S3UploadClientHandler{
		S3UploadClientUseCase: sUseCase,
		StorageBucketUseCase:  storageUC,
	}
	v1 := e.Group("/api/v1")
	{
		s3Upload := v1.Group("/s3Upload")
		{
			s3Upload.GET("/list", handler.ListObjects)
			s3Upload.POST("/detail", handler.Detail)
			s3Upload.POST("/rename", handler.ReName)
			s3Upload.POST("/remove", handler.Delete)
			s3Upload.POST("/upload", handler.Upload)
			s3Upload.GET("/download", handler.Download)
		}
	}
}

// ListObjects @Title ListObjects
// @Description 查询文件列表
// @Accept  json
// @Tags  s3Upload
// @Param bucketID formData integer true "储存桶id"
// @Param filePrefix formData string true "前缀"
// @Success 200 {object} []domain.UploadedFile "查询成功"
// @Failure 400	"查询失败"
// @Router /api/v1/s3Upload/list  [GET]
func (s *S3UploadClientHandler) ListObjects(c echo.Context) error {
	bucketID, _ := strconv.Atoi(c.QueryParam("bucketID"))
	filePrefix := c.QueryParam("filePrefix")
	ctx := c.Request().Context()
	storage, err := s.StorageBucketUseCase.GetByID(ctx, uint(bucketID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.S3UploadCode, err))
	}
	var s3Upload domain.S3Upload
	err = json.Unmarshal([]byte(storage.Config), &s3Upload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.S3UploadCode, domain.JsonUnMarshalErr))
	}
	fileNames, err := s.S3UploadClientUseCase.S3UploadFiles(ctx, s3Upload, filePrefix)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.S3UploadCode, err))
	}
	result := make([]domain.UploadedFile, 0)
	for _, name := range fileNames {
		file, err := s.S3UploadClientUseCase.Detail(ctx, s3Upload, name)
		if err != nil {
			//return c.JSON(http.StatusBadRequest, GetResponseErr(domain.S3UploadCode, err))
		}
		result = append(result, file)
	}
	return c.JSON(http.StatusOK, SuccessWriteResult(result))
}

// ReName @Title ReName
// @Description 重命名
// @Accept  json
// @Tags  s3Upload
// @Param bucketID formData integer true "储存桶id"
// @Param oldName formData string true "旧名称"
// @Param newName formData string true "新名词"
// @Success 200 "重命名成功"
// @Failure 400	"重命名失败"
// @Router /api/v1/s3Upload/rename  [PUT]
func (s *S3UploadClientHandler) ReName(c echo.Context) error {
	param := struct {
		BucketID uint   `json:"bucketID"`
		OldName  string `json:"oldName"`
		NewName  string `json:"newName"`
	}{}
	err := c.Bind(&param)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.S3UploadCode, domain.ParamErr))
	}
	ctx := c.Request().Context()
	storage, err := s.StorageBucketUseCase.GetByID(ctx, param.BucketID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.S3UploadCode, err))
	}

	var s3Upload domain.S3Upload
	err = json.Unmarshal([]byte(storage.Config), &s3Upload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.S3UploadCode, domain.JsonUnMarshalErr))
	}
	err = s.S3UploadClientUseCase.ReName(ctx, s3Upload, param.OldName, param.NewName)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.S3UploadCode, err))
	}
	return c.JSON(http.StatusOK, SuccessResult())
}

// Delete @Title Delete
// @Description 删除
// @Accept  json
// @Tags  s3Upload
// @Param bucketID formData integer true "储存桶id"
// @Param fileName formData string true "文件名称"
// @Success 200 "删除成功"
// @Failure 400	"删除失败"
// @Router /api/v1/s3Upload/remove  [POST]
func (s *S3UploadClientHandler) Delete(c echo.Context) error {
	param := struct {
		BucketID uint   `json:"bucketID"`
		FileName string `json:"fileName"`
	}{}
	err := c.Bind(&param)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.S3UploadCode, domain.ParamErr))
	}
	ctx := c.Request().Context()
	storage, err := s.StorageBucketUseCase.GetByID(ctx, param.BucketID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.S3UploadCode, err))
	}

	var s3Upload domain.S3Upload
	err = json.Unmarshal([]byte(storage.Config), &s3Upload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.S3UploadCode, domain.JsonUnMarshalErr))
	}
	err = s.S3UploadClientUseCase.Delete(ctx, s3Upload, param.FileName)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.S3UploadCode, err))
	}
	return c.JSON(http.StatusOK, SuccessResult())
}

// Upload @Title Upload
// @Description oss上传
// @Accept  json
// @Tags  s3Upload
// @Param bucketID formData integer true "储存桶id"
// @Param path formData string true "oss路径"
// Param file formData string true "文件"
// @Success 200 "上传成功"
// @Failure 400	"上传失败"
// @Router /api/v1/s3Upload/upload  [POST]
func (s *S3UploadClientHandler) Upload(c echo.Context) error {
	bucketID, _ := strconv.Atoi(c.FormValue("bucketID"))
	path := c.FormValue("path")
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.S3UploadCode, domain.ParamErr))
	}
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.S3UploadCode, domain.FileOpenErr))
	}

	writePath := fmt.Sprintf("%s/%s", utils.GetOSSUploadPath(), file.Filename)
	dst, err := os.Create(writePath)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.S3UploadCode, domain.FileCreateErr))
	}

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.S3UploadCode, domain.FileCopyErr))
	}

	src.Close()
	dst.Close()
	defer func() {
		os.Remove(writePath)
	}()

	ctx := c.Request().Context()
	storage, err := s.StorageBucketUseCase.GetByID(ctx, uint(bucketID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.S3UploadCode, err))
	}

	var s3Upload domain.S3Upload
	err = json.Unmarshal([]byte(storage.Config), &s3Upload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.S3UploadCode, domain.JsonUnMarshalErr))
	}
	err = s.S3UploadClientUseCase.Upload(ctx, s3Upload, path, file.Filename)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.S3UploadCode, err))
	}

	return c.JSON(http.StatusOK, SuccessResult())
}

// Download @Title Download
// @Description 下载
// @Accept  json
// @Tags  s3Upload
// @Param bucketID formData integer true "储存桶id"
// @Param fileName formData string true "文件名称"
// @Success 200 "下载成功"
// @Failure 400	"下载失败"
// @Router /api/v1/s3Upload/download  [GET]
func (s *S3UploadClientHandler) Download(c echo.Context) error {
	param := struct {
		BucketID uint   `json:"bucketID"`
		FileName string `json:"fileName"`
	}{}
	err := c.Bind(&param)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.S3UploadCode, domain.ParamErr))
	}
	ctx := c.Request().Context()
	storage, err := s.StorageBucketUseCase.GetByID(ctx, param.BucketID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.S3UploadCode, err))
	}

	var s3Upload domain.S3Upload
	err = json.Unmarshal([]byte(storage.Config), &s3Upload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.S3UploadCode, domain.JsonUnMarshalErr))
	}
	fileInfo, err := s.S3UploadClientUseCase.Download(ctx, s3Upload, param.FileName)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.S3UploadCode, err))
	}
	return c.JSON(http.StatusOK, SuccessWriteResult(fileInfo))
}

// Detail @Title Detail
// @Description 文件详情
// @Accept  json
// @Tags  s3Upload
// @Param bucketID formData integer true "储存桶id"
// @Param fileName formData string true "文件名称"
// @Success 200 {object} domain.UploadedFile "查询成功"
// @Failure 400	"查询失败"
// @Router /api/v1/s3Upload/detail  [POST]
func (s *S3UploadClientHandler) Detail(c echo.Context) error {
	param := struct {
		BucketID uint   `json:"bucketID"`
		FileName string `json:"fileName"`
	}{}
	err := c.Bind(&param)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.S3UploadCode, domain.ParamErr))
	}

	ctx := c.Request().Context()
	storage, err := s.StorageBucketUseCase.GetByID(ctx, param.BucketID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.S3UploadCode, err))
	}

	var s3Upload domain.S3Upload
	err = json.Unmarshal([]byte(storage.Config), &s3Upload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.S3UploadCode, domain.JsonUnMarshalErr))
	}
	fileInfo, err := s.S3UploadClientUseCase.Detail(ctx, s3Upload, param.FileName)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GetResponseErr(domain.S3UploadCode, err))
	}
	return c.JSON(http.StatusOK, SuccessWriteResult(fileInfo))
}
