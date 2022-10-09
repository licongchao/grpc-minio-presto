package httpserver

import (
	dlsvc "da/datalakesvc"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func GetHTTPServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/upload", func(w http.ResponseWriter, r *http.Request) {
		uuid, _ := dlsvc.GrpcSvc.PrepareRawData(r)

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(uuid); err != nil {
			fmt.Print(err)
		}
	})
	mux.HandleFunc("/api/query/", func(w http.ResponseWriter, r *http.Request) {
		pathSplit := strings.Split(r.URL.Path, "/")
		queryResp, _ := dlsvc.GrpcSvc.GetDataFromUUID(pathSplit[len(pathSplit)-1])

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		w.Write([]byte(queryResp))
		// if err := json.NewEncoder(w).Encode(queryResp); err != nil {
		// 	fmt.Print(err)
		// }
	})

	return mux
}
