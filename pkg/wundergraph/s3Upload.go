package wundergraph

type S3UploadProvider struct {
	Name            string `json:"name"`
	Endpoint        string `json:"endpoint"`
	AccessKeyID     string `json:"accessKeyID"`
	SecretAccessKey string `json:"secretAccessKey"`
	BucketName      string `json:"bucketName"`
	BucketLocation  string `json:"bucketLocation"`
	UseSSL          bool   `json:"useSSL"`
}
