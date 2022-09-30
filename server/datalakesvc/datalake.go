package datalakesvc

import (
	"encoding/json"
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
		表结构定义分为如下两种
		-- map[string]interface {} / []interface {}
		map类型需要传入完整的表结构定义
			示例如下:
			{
			  	sepal_length DOUBLE,
				sepal_width  DOUBLE,
				petal_length DOUBLE,
				petal_width  DOUBLE,
				class        VARCHAR
			}
		[] 类型将默认采用String实现
	(3)写入UUID.json
	返回UUID给后续使用
*/
func (s *DatalakeGRPCSvc) PrepareRawData(r *http.Request) (uuidStr string, err error) {
	uuidExchange, _ := ObjStorageSvc.UploadObject(r)
	newUUID := uuid.New()

	keyName := newUUID.String() + ".json"
	uuidJson, _ := json.Marshal(uuidExchange)
	ObjStorageSvc.CreateObject(uuidJson, BucketName, keyName)

	return newUUID.String(), nil
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
