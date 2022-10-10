package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}
func Jsonify(rows *sql.Rows) (string, error) {
	columns, err := rows.Columns()
	if err != nil {
		return "", err
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	jsonData, err := json.Marshal(tableData)
	if err != nil {
		return "", err
	}
	fmt.Println(string(jsonData))
	return string(jsonData), nil
}

// func Jsonify(rows *sql.Rows) []string {
// 	columns, err := rows.Columns()
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	values := make([]interface{}, len(columns))

// 	scanArgs := make([]interface{}, len(values))
// 	for i := range values {
// 		scanArgs[i] = &values[i]
// 	}

// 	c := 0
// 	results := make(map[string]interface{})
// 	data := []string{}

// 	for rows.Next() {
// 		if c > 0 {
// 			data = append(data, ",")
// 		}

// 		err = rows.Scan(scanArgs...)
// 		if err != nil {
// 			panic(err.Error())
// 		}

// 		for i, value := range values {
// 			switch value.(type) {
// 			case nil:
// 				results[columns[i]] = nil

// 			case []byte:
// 				s := string(value.([]byte))
// 				x, err := strconv.Atoi(s)

// 				if err != nil {
// 					results[columns[i]] = s
// 				} else {
// 					results[columns[i]] = x
// 				}

// 			default:
// 				results[columns[i]] = value
// 			}
// 		}

// 		b, _ := json.Marshal(results)
// 		data = append(data, strings.TrimSpace(string(b)))
// 		c++
// 	}

// 	return data
// }
