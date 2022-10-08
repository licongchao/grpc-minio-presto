package main

import (
	dlsvc "da/datalakesvc"
	mylog "da/mylog"
	"os"

	"github.com/urfave/cli"
)

var (
	dsn      = "http://minio_access_key@localhost:8888?catalog=datalake"
	S3config = &dlsvc.SvcConnectionConfig{Endpoint: "http://127.0.0.1:9000", Region: "dai-cn", Access_key_id: "minio_access_key", Access_key_secret: "minio_secret_key"}
)

// func GetHTTPServeMux() *http.ServeMux {
// 	mux := http.NewServeMux()
// 	mux.HandleFunc("/api/upload", func(w http.ResponseWriter, r *http.Request) {
// 		out, _ := dlsvc.ConnSvc.ExePrestoSqlQuery("select * from postgresql.public.iris_meta")
// 		fmt.Println(string(out))
// 	})

// 	return mux
// }

func main() {
	dlsvc.DatalakePrefix = ""
	dlsvc.BucketName = "datalake"

	dlsvc.ConnSvc, _ = dlsvc.InitConnection(dsn)
	dlsvc.ObjStorageSvc, _ = dlsvc.ConnectStorageSvc(*S3config)
	mylog.InitLogger()

	// s3svc.ObjStorageSvc.ListObjectNames(&s3.ListObjectsInput{
	// 	Bucket: aws.String("datalake"),
	// 	Prefix: aws.String("rawdata"),
	// })

	app := cli.NewApp()
	app.Commands = []cli.Command{
		StartServerCommand(),
	}
	app.Run(os.Args)

	// dsn := "http://minio_access_key@localhost:8888?catalog=postgresql&schema=public"
	// sqlExe := "insert into postgresql.public.iris_meta(name,price) values('hello', 0.2)"
	// dsn := "http://minio_access_key@localhost:8888?catalog=datalake&schema=iris"
	// // sqlExe := "insert into datalake.iris.iris_parquet(name,price) values('hello', 0.2)"
	// // dsn := "http://dummy@localhost:8888"
	// sqlExe := "select * from datalake.iris.iris_parquet"
	// sqlExe := "insert into datalake.iris.iris_parquet(class,petal_length,petal_width,sepal_length,sepal_width) values('Hello', 0.1,0.2,0.3,0.4)"
	// // sqlExe := "select * from postgresql.public.iris_meta"
	// sqlExe := "select * from datalake.iris.iris_parquet l left join postgresql.public.iris_meta r on l.class=r.name where name is not null"
}
