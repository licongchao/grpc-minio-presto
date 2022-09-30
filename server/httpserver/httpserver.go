package httpserver

import (
	dlsvc "da/datalakesvc"
	"encoding/json"
	"net/http"
)

func GetHTTPServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/upload", func(w http.ResponseWriter, r *http.Request) {
		uuid, _ := dlsvc.GrpcSvc.PrepareRawData(r)

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(uuid); err != nil {
			panic(err)
		}
	})

	return mux
}
