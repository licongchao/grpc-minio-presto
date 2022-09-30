package datalakesvc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
)

var (
	S3config       *SvcConnectionConfig
	ObjStorageSvc  ObjectStorageGRPCSvc
	BucketName     string
	DatalakePrefix string
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
	Alias       string //在Datalake中的资源别名Table
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

	resp, err := s.s3Client.ListBuckets(params)
	if err != nil {
		err = errors.Wrapf(err, "ListBuckets Failed")
		return buckets, err
	}

	for _, bucket := range resp.Buckets {
		buckets = append(buckets, *bucket.Name)
	}
	return buckets, nil
}

// Object Handy Tools
func (s *ObjectStorageGRPCSvc) ListObjectNames(params *s3.ListObjectsInput) (ObjectNames []string, err error) {
	resp, err := s.s3Client.ListObjects(params)
	if err != nil {
		err = errors.Wrapf(err, "ListObjects Failed")
		return ObjectNames, err
	}
	for _, obj := range resp.Contents {
		ObjectNames = append(ObjectNames, *obj.Key)
	}
	return ObjectNames, nil
}

func (s *ObjectStorageGRPCSvc) CreateObject(fileContent []byte, bucketName string, keyName string) (err error) {
	s.s3Client.PutObject(&s3.PutObjectInput{
		Body:   bytes.NewReader(fileContent),
		Bucket: aws.String(bucketName),
		Key:    aws.String(keyName),
	})
	return err
}

/*
	为了简化程序设计和处理的复杂度，
	文件名以Alias.EXTENSION 定义
	保存在对象存储的文件名 = Alias_Timestamp

	如果需要自定义Schema，需要在http头中定义
	{
		"sepal_length": "DOUBLE",
		"sepal_width":  "DOUBLE",
		"petal_length": DOUBLE,
		"petal_width":  DOUBLE,
		"class":        VARCHAR
	},
	如果是Array类型， 则默认初始化为VARCHAR
*/
func createSchemaColumns(schema string) (columns string) {
	var column_list = []string{}
	var jsonObj interface{}
	json.Unmarshal([]byte(schema), &jsonObj)

	if jsonObj == nil {
		log.Print("Not supportted Schema")
	} else if reflect.TypeOf(jsonObj).String() == "map[string]interface {}" {
		var result map[string]string
		json.Unmarshal([]byte(schema), &result)

		for k, v := range result {
			column_list = append(column_list, strings.ReplaceAll(k, "\"", "")+" "+strings.ReplaceAll(v, "\"", ""))
		}
	} else if reflect.TypeOf(jsonObj).String() == "[]interface {}" {
		var result []string
		json.Unmarshal([]byte(schema), &result)
		for _, v := range result {
			column_list = append(column_list, strings.ReplaceAll(v, "\"", "")+" VARCHAR")
		}
	} else {
		log.Print("Not supportted Schema")
	}
	return strings.Join(column_list, ",")
}

func createSchemaSTMT(schema string, filename string) (schemaSTMT string) {
	var stmt strings.Builder

	if schema == "" {
		log.Print("Cannot find Schema definition, Will skip create Schema Process")
		return ""
	} else {
		alias := strings.Split(filename, ".")[0]
		suffix := strings.Split(filename, ".")[1]
		tableName := alias + "_" + suffix
		stmt.WriteString(fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS datalake.%s WITH (location = 's3a://datalake/');`, alias))
		stmt.WriteString(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS datalake.%s.%s (%s)`, alias, tableName, createSchemaColumns(schema)))
		stmt.WriteString(fmt.Sprintf(`WITH (external_location = 's3a://datalake/%s/%s',format = '%s');`, alias, tableName, suffix))
	}
	return stmt.String()
}

func (s *ObjectStorageGRPCSvc) UploadObject(r *http.Request) (uuidExchange UUIDExchange, err error) {
	var buf bytes.Buffer
	r.ParseMultipartForm(32 << 20)

	file, header, err := r.FormFile("file")
	// Upload file first
	fileNameSplits := strings.Split(header.Filename, ".")
	newFilename := fmt.Sprintf("%s_%d.%s", fileNameSplits[0], time.Now().Unix(), fileNameSplits[1])
	io.Copy(&buf, file)
	_, errPutObject := s.s3Client.PutObject(&s3.PutObjectInput{
		Body:   bytes.NewReader(buf.Bytes()),
		Bucket: aws.String(BucketName),
		Key:    aws.String(DatalakePrefix + fileNameSplits[0] + "/" + fileNameSplits[0] + "_" + fileNameSplits[1] + "/" + newFilename),
	})
	// Create Mapping Schema
	schema := r.FormValue("schema")
	stmt := createSchemaSTMT(schema, header.Filename)
	log.Print(stmt)
	ConnSvc.ExePrestoSqlQuery(stmt)

	if err != nil {
		return UUIDExchange{}, err
	}
	defer file.Close()
	return UUIDExchange{
			Filename:    header.Filename,
			DatalakeSrc: DatalakePrefix + newFilename,
			Alias:       fileNameSplits[0],
			DAG:         "",
			Sql:         "",
		},
		errPutObject
}