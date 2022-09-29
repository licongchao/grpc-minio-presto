package datalakesvc

import (
	"io"
	"net/http"

	"github.com/google/uuid"
	_ "github.com/prestodb/presto-go-client/presto"
)

type UUIDMetaInfo struct {
	Datasource string
	Airflow    string
	Sql        string
}
type DatalakeGRPCSvc struct {
}

/*
	high level API
	(1) 获取datalake/UUID.json
	(2) 解析UUID.json内容
	{
		datasource: datalake/rawdata/xxxxxx.csv
		airflow: DAG-ID
		sql: select * from xxxx
	}
	(3) 调用数据库查询条件并返回String
*/
// func (s *DatalakeGRPCSvc) GetDataFromUUID(uuid string) (data string, err error) {

// }

/*
	high level API
	(1)上传文件
	(2)自动创建表结构
	(3)写入UUID.json
	返回UUID给后续使用
*/
func (s *DatalakeGRPCSvc) PrepareRawData(r *http.Request) (uuidStr string, err error) {
	r.ParseMultipartForm(32 << 20)

	file, handler, err := r.FormFile("file")

	if err != nil {
		return "", err
	}
	defer file.Close()

	io.Copy()
	newUUID := uuid.New()

	UUIDMetaInfo{"datalake/rawdata/" + handler.Filename, "", tableAlias}

	return
}

// low level API,直接调用SQL查询所需的数据
// func (s *DatalakeGRPCSvc) GetDataViaSql(sqlStr string) (data string, err error) {
// }

// func getCsvColumns(r *io.Reader) ([]string, error) {
// 	var cols []string

// 	rec, err := r.Read()

// 	if err != nil {
// 		return cols, err
// 	}
// 	return rec, nil
// }

func main() {
	// f, err := os.Open("/home/licongchao/workspace/DataPreparation/presto-minio-docker/data/titanicWithCols.csv")
	// if err != nil {
	// 	return
	// }
	// csvReader := csv.NewReader(f)
	// // for {
	// // rec, err := csvReader.Read()
	// // if err == io.EOF {
	// // 	return
	// // }
	// // if err != nil {
	// // 	log.Fatal(err)
	// // }
	// // // do something with read line
	// // fmt.Printf("%+v\n", rec)
	// // }

	// cols, err := getCsvColumns(csvReader)
	// fmt.Println(cols)
}
