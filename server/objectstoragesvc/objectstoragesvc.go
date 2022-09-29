package objectstoragesvc

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
)

type SvcConnectionConfig struct {
	Endpoint          string
	Region            string
	Access_key_id     string
	Access_key_secret string
}
type ObjectStorageGRPCSvc struct {
	s3Client *s3.S3
	endpoint string
}

// UUID设计，将请求保留
type UUIDExchange struct {
	Filename    string //用户上传的原始文件名
	DatalakeSrc string //在Datalake中的路径
	DAG         string //如果需要做ETL，则存放DAG ID
	Sql         string //Sql语句
}

/**
	s3config := s3svc.SvcConnectionConfig{Endpoint: "http://127.0.0.1:9000", Region: "dai-cn", Access_key_id: "minio_access_key", Access_key_secret: "minio_secret_key"}
	s3connection, _ := s3svc.ConnectStorageSvc(s3config)
**/
// ObjectStorage Service Connection
func ConnectStorageSvc(config SvcConnectionConfig) (svc ObjectStorageGRPCSvc, err error) {
	creds := credentials.NewStaticCredentials(config.Access_key_id, config.Access_key_secret, "")
	awsconfig := &aws.Config{
		Endpoint:         &config.Endpoint,
		Region:           aws.String(config.Region),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      creds,
		DisableSSL:       aws.Bool(true),
	}
	newSession, err := session.NewSession(awsconfig)
	if err != nil {
		err = errors.Wrapf(err,
			"Failed to CreateSession on %s",
			config.Endpoint)
		return
	}
	s3Client := s3.New(newSession)

	svc.endpoint = config.Endpoint
	svc.s3Client = s3Client
	return
}

/*
    Bucket Handy Tools

	s3connection.CreateBucket(&s3.CreateBucketInput{Bucket: aws.String("helloworld")})
*/
func (s *ObjectStorageGRPCSvc) CreateBucket(params *s3.CreateBucketInput) (err error) {
	_, err = s.s3Client.CreateBucket(params)
	if err != nil {
		err = errors.Wrapf(err, "CreateBucket Failed")
	}
	return
}

/*
	s3connection.CreateBucket(&s3.CreateBucketInput{Bucket: aws.String("helloworld")})
*/
func (s *ObjectStorageGRPCSvc) DeleteBucket(params *s3.DeleteBucketInput) (err error) {
	_, err = s.s3Client.DeleteBucket(params)
	if err != nil {
		err = errors.Wrapf(err, "DeleteBucket Failed")
	}
	return
}

/*
	buckets, err := s3connection.ListBuckets(&s3.ListBucketsInput{})
*/
func (s *ObjectStorageGRPCSvc) ListBuckets(params *s3.ListBucketsInput) (buckets []string, err error) {
	var bucketNames = []string{}

	resp, err := s.s3Client.ListBuckets(params)
	if err != nil {
		err = errors.Wrapf(err, "ListBuckets Failed")
		return bucketNames, err
	}

	for _, bucket := range resp.Buckets {
		bucketNames = append(bucketNames, *bucket.Name)
	}
	return bucketNames, nil
}

// Object Handy Tools
func (s *ObjectStorageGRPCSvc) ListObjects(params *s3.ListObjectsInput) (err error) {
	_, err = s.s3Client.ListObjects(params)
	if err != nil {
		err = errors.Wrapf(err, "ListObjects Failed")
	}
	return
}

/*
	为了简化程序设计和处理的复杂度，将上传的文件名以如下规则定义
	Alias_Timestamp
*/
var (
	bucketName     = "datalake"
	datalakePrefix = "rawdata/"
)

func (s *ObjectStorageGRPCSvc) UploadObject(r *http.Request) (err error) {
	var buf bytes.Buffer
	r.ParseMultipartForm(32 << 20)

	file, header, err := r.FormFile("file")
	if err != nil {
		return err
	}
	defer file.Close()

	fileNameSplits := strings.Split(header.Filename, ".")
	newFilename := fmt.Sprintf("%s_%d.%s", fileNameSplits[0], time.Now().Unix(), fileNameSplits[1])
	io.Copy(&buf, file)
	s.s3Client.PutObject(&s3.PutObjectInput{
		Body:   bytes.NewReader(buf.Bytes()),
		Bucket: aws.String(bucketName),
		Key:    aws.String(datalakePrefix + newFilename),
	})
	return nil
}
