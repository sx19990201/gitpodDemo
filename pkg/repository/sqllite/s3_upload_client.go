package sqllite

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/fire_boom/domain"
	"github.com/fire_boom/utils"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	log "github.com/sirupsen/logrus"
	"io"
	"strings"
	"sync"
)

const (
	DateFormatStr = "2006-01-02 15:04:05"
)

var minioClientMap sync.Map

type S3UploadClientRepository struct {
}

func NewS3UploadClientRepository() *S3UploadClientRepository {
	return &S3UploadClientRepository{}
}

func (s *S3UploadClientRepository) S3UploadFiles(ctx context.Context, s3Upload domain.S3Upload, fileName string) (result []string, err error) {
	client, err := NewMinioClient(s3Upload)
	if err != nil {
		log.Error("S3UploadClientRepository S3UploadFiles NewMinioClient err : ", err.Error())
		return
	}
	options := minio.ListObjectsOptions{
		Prefix:    fileName,
		Recursive: false,
	}
	for info := range client.ListObjects(context.TODO(), s3Upload.BucketName, options) {
		result = append(result, info.Key)
	}
	return
}

func (s *S3UploadClientRepository) Copy(ctx context.Context, s3Upload domain.S3Upload, oldName, newName string) (err error) {
	client, err := NewMinioClient(s3Upload)
	if err != nil {
		log.Error("S3UploadClientRepository Copy NewMinioClient err : ", err.Error())
		return
	}
	// 设置源对象文件
	src := minio.CopySrcOptions{
		Bucket: s3Upload.BucketName,
		Object: oldName,
		Start:  0,
		End:    20,
	}
	// 设置目标文件
	dest := minio.CopyDestOptions{
		Bucket: s3Upload.BucketName,
		Object: newName,
	}
	_, err = client.ComposeObject(ctx, dest, src)
	// 内容copy
	//_, err = client.CopyObject(ctx, dest, src)
	if err != nil {
		log.Error("S3UploadClientRepository Copy ComposeObject err : ", err.Error())
		return
	}
	return
}

func (s *S3UploadClientRepository) CreateFile(ctx context.Context, s3Upload domain.S3Upload, ossPath, fileName string) (err error) {
	client, err := NewMinioClient(s3Upload)
	if err != nil {
		log.Error("S3UploadClientRepository CreateFile NewMinioClient err : ", err.Error())
		return
	}
	objectName := fmt.Sprintf("%s%s", ossPath, fileName)
	filePath := fmt.Sprintf("%s/%s", utils.GetOSSUploadPath(), fileName)
	// 创建文件
	_, err = client.FPutObject(ctx, s3Upload.BucketName, objectName, filePath, minio.PutObjectOptions{ContentDisposition: "inline"})
	if err != nil {
		log.Error("S3UploadClientRepository CreateFile FPutObject err : ", err.Error())
		return
	}
	return
}

func (s *S3UploadClientRepository) DeleteFiles(ctx context.Context, s3Upload domain.S3Upload, fileName string) (err error) {
	client, err := NewMinioClient(s3Upload)
	if err != nil {
		log.Error("S3UploadClientRepository DeleteFiles FPutObject err : ", err.Error())
		return
	}
	objectsCh := make(chan minio.ObjectInfo)
	options := minio.ListObjectsOptions{
		Prefix:    fileName,
		Recursive: true,
	}
	go func() {
		defer close(objectsCh)
		for object := range client.ListObjects(ctx, s3Upload.BucketName, options) {
			if object.Err != nil {
				log.Error("S3UploadClientRepository DeleteFiles FPutObject err : ", object.Err)
			}
			objectsCh <- object
		}
	}()

	for rErr := range client.RemoveObjects(ctx, s3Upload.BucketName, objectsCh, minio.RemoveObjectsOptions{}) {
		log.Error("Error detected during deletion: ", rErr)
	}
	return
}

func (s *S3UploadClientRepository) Download(ctx context.Context, s3Upload domain.S3Upload, filePath string) (result []string, err error) {
	client, err := NewMinioClient(s3Upload)
	if err != nil {
		log.Error("S3UploadClientRepository Download NewMinioClient err : ", err.Error())
		return nil, err
	}
	file, err := client.GetObject(ctx, s3Upload.BucketName, filePath, minio.GetObjectOptions{})
	if err != nil {
		log.Error("S3UploadClientRepository Download GetObject err : ", err.Error())
		return nil, err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		data, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return result, err
			}
		}
		result = append(result, string(data))
	}
	return
}

func (s *S3UploadClientRepository) Detail(ctx context.Context, s3Upload domain.S3Upload, fileName string) (result domain.UploadedFile, err error) {
	client, err := NewMinioClient(s3Upload)
	if err != nil {
		log.Error("S3UploadClientRepository Detail NewMinioClient err : ", err.Error())
		return
	}
	fileInfo, err := client.StatObject(ctx, s3Upload.BucketName, fileName, minio.StatObjectOptions{})
	if err != nil {
		log.Error("S3UploadClientRepository Detail StatObject err : ", err.Error())
		return
	}
	result.Name = fileInfo.Key
	result.Size = fileInfo.Size
	result.CreateTime = fileInfo.LastModified.Format(DateFormatStr)
	result.UpdateTime = fileInfo.Expiration.Format(DateFormatStr)
	result.URL = fmt.Sprintf("https://%s.%s/%s", s3Upload.BucketName, s3Upload.EndPoint, fileName)
	result.MimeTypes = fileInfo.ContentType
	arr := strings.Split(fileInfo.Key, "")
	if arr[len(arr)-1] == "/" {
		result.IsDir = true
	}
	//if strings.Contains(fileInfo.Key, ".jpg") {
	//
	//}
	//if strings.Contains(fileInfo.Key, ".jpeg") {
	//	result.MimeTypes = "image/jpg"
	//}
	//if strings.Contains(fileInfo.Key, ".gif") {
	//	result.MimeTypes = "image/gif"
	//}
	//if strings.Contains(fileInfo.Key, ".png") {
	//	result.MimeTypes = "image/png"
	//}

	return
}

func NewMinioClient(s3Upload domain.S3Upload) (client *minio.Client, err error) {
	if val, ok := minioClientMap.Load(s3Upload.Name); ok {
		client, ok = val.(*minio.Client)
		if !ok {
			return nil, errors.New("get client fail")
		}
		return
	}
	client, err = minio.New(s3Upload.EndPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(s3Upload.AccessKeyID, s3Upload.SecretAccessKey, ""),
		Secure: s3Upload.UseSSL,
	})
	minioClientMap.Store(s3Upload.Name, client)
	return
}
