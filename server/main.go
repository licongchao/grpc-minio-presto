package main

// import (
// 	"database/sql"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"os"
// 	"reflect"

// 	s3svc "da/objectstoragesvc"

// 	"github.com/go-sql-driver/mysql"
// 	_ "github.com/prestodb/presto-go-client/presto"
// 	"github.com/urfave/cli"
// )
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	s3svc "da/objectstoragesvc"

	"github.com/go-sql-driver/mysql"
	_ "github.com/prestodb/presto-go-client/presto"
)

var (
	s3config        = s3svc.SvcConnectionConfig{Endpoint: "http://127.0.0.1:9000", Region: "dai-cn", Access_key_id: "minio_access_key", Access_key_secret: "minio_secret_key"}
	s3connection, _ = s3svc.ConnectStorageSvc(s3config)
)

func GetHTTPServeMux() *http.ServeMux {

	mux := http.NewServeMux()
	mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		s3connection.UploadObject(r)
	})

	return mux
}

func main() {
	// app := cli.NewApp()
	// app.Commands = []cli.Command{
	// 	StartServerCommand(),
	// }
	// app.Run(os.Args)

	// dsn := "http://minio_access_key@localhost:8888?catalog=postgresql&schema=public"
	// sqlExe := "insert into postgresql.public.iris_meta(name,price) values('hello', 0.2)"
	dsn := "http://minio_access_key@localhost:8888?catalog=datalake&schema=iris"
	// sqlExe := "insert into datalake.iris.iris_parquet(name,price) values('hello', 0.2)"
	// dsn := "http://dummy@localhost:8888"
	sqlExe := "select * from datalake.iris.iris_parquet"
	// sqlExe := "insert into datalake.iris.iris_parquet(class,petal_length,petal_width,sepal_length,sepal_width) values('Hello', 0.1,0.2,0.3,0.4)"
	// // sqlExe := "select * from postgresql.public.iris_meta"
	// sqlExe := "select * from datalake.iris.iris_parquet l left join postgresql.public.iris_meta r on l.class=r.name where name is not null"

	// db, err := sql.Open("presto", dsn)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// rows, err := db.Query(sqlExe)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// var (
	// 	sepal_length float64
	// 	sepal_width  float64
	// 	petal_length float64
	// 	petal_width  float64
	// 	class        string
	// )
	// defer rows.Close()
	// for rows.Next() {
	// 	log.Println("hello")
	// 	err := rows.Scan(&sepal_length, &sepal_width, &petal_length, &petal_width, &class)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Println(sepal_length, sepal_width, petal_length, petal_width, class)
	// }
	// err = rows.Err()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// defer db.Close()

	content, _ := ExePrestoSqlQuery(dsn, sqlExe)
	fmt.Println("query result :  ", string(content))
}

type jsonNullInt64 struct {
	sql.NullInt64
}

func (v jsonNullInt64) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(v.Int64)
}

type jsonNullFloat64 struct {
	sql.NullFloat64
}

func (v jsonNullFloat64) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(v.Float64)
}

type jsonNullTime struct {
	mysql.NullTime
}

func (v jsonNullTime) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(v.Time)
}

// --------------------------------------------------------------
var jsonNullInt64Type = reflect.TypeOf(jsonNullInt64{})
var jsonNullFloat64Type = reflect.TypeOf(jsonNullFloat64{})
var jsonNullTimeType = reflect.TypeOf(jsonNullTime{})
var nullInt64Type = reflect.TypeOf(sql.NullInt64{})
var nullFloat64Type = reflect.TypeOf(sql.NullFloat64{})
var nullTimeType = reflect.TypeOf(mysql.NullTime{})

func ExePrestoSqlQuery(prestoUrl string, sqlExe string) ([]byte, error) {
	db, err := sql.Open("presto", prestoUrl)
	if err != nil {
		return nil, fmt.Errorf("can't connect to presto error: %v", err)
	}
	rows, err := db.Query(sqlExe)
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("column error: %v", err)
	}

	ct, err := rows.ColumnTypes()
	if err != nil {
		return nil, fmt.Errorf("column type error: %v", err)
	}

	types := make([]reflect.Type, len(ct))
	for i, tp := range ct {
		st := tp.ScanType()
		if st == nil {
			return nil, fmt.Errorf("scantype is null for column: %v", err)
		}
		switch st {
		case nullInt64Type:
			types[i] = jsonNullInt64Type
		case nullFloat64Type:
			types[i] = jsonNullFloat64Type
		case nullTimeType:
			types[i] = jsonNullTimeType
		default:
			types[i] = st
		}
	}
	values := make([]interface{}, len(ct))
	var slice []map[string]interface{}
	for rows.Next() {
		for i := range values {
			values[i] = reflect.New(types[i]).Interface()
		}
		err = rows.Scan(values...)
		if err != nil {
			return nil, fmt.Errorf("failed to scan values: %v", err)
		}
		data := make(map[string]interface{})
		for i, v := range values {
			data[columns[i]] = v
		}
		slice = append(slice, data)
	}

	return json.Marshal(slice)
}
