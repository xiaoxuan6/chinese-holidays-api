package router

import (
	"github.com/gorilla/mux"
	"github.com/xiaoxuan6/chinese-holidays-api/handles"
	"net/http"
)

func InitRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(mux.CORSMethodMiddleware(r))

	r.NotFoundHandler = http.HandlerFunc(handles.NotFoundHandle)

	r.HandleFunc("/", handles.IndexHandle).Methods(http.MethodGet)

	r.HandleFunc("/api/holidays", handles.DateHandle).Methods(http.MethodGet).Queries("date", "{date}")
	r.HandleFunc("/api/holidays/{date}", handles.DateHandle).Methods(http.MethodGet)

	return r
}
