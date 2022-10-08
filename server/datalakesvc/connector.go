package datalakesvc

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/go-sql-driver/mysql"
)

var (
	ConnSvc DBConnectorSvc
	GrpcSvc DatalakeGRPCSvc
)

type jsonNullInt64 struct {
	sql.NullInt64
}

type DBConnectorSvc struct {
	conn *sql.DB
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

func InitConnection(url string) (DBConnectorSvc, error) {
	db, err := sql.Open("trino", url)
	if err != nil {
		newErr := fmt.Errorf("can't connect to presto error: %v", err)
		log.Print(newErr)
		return DBConnectorSvc{nil}, newErr
	}
	ConnSvc.conn = db
	return ConnSvc, nil
}

func CloseRow(rows *sql.Rows) {
	if rows != nil {
		rows.Close()
	}
}

// func (s *DBConnectorSvc) ExecPrestoSql(sqlExe string) error {
// 	result, err := ConnSvc.conn.Exec(sqlExe)
// 	if err != nil {
// 		fmt.Print(err)
// 		return err
// 	}
// 	fmt.Print(result)
// 	return nil
// }

func (s *DBConnectorSvc) ExecPrestoSqlQuery(sqlExe string) ([]byte, error) {
	rows, err := ConnSvc.conn.Query(sqlExe)
	if err != nil {
		fmt.Print(err)
	}
	defer CloseRow(rows)

	if rows == nil {
		return []byte{}, nil
	}
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
