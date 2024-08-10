package router

import (
	"github.com/gorilla/mux"
	"github.com/xiaoxuan6/chinese-holidays-api/handlers"
	"net/http"
)

func InitRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(mux.CORSMethodMiddleware(r))

	r.NotFoundHandler = http.HandlerFunc(handlers.NotFoundHandler)

	r.HandleFunc("/", handlers.IndexHandler).Methods(http.MethodGet)

	r.HandleFunc("/api/holidays", handlers.DateHandler).Methods(http.MethodGet).Queries("date", "{date}")
	r.HandleFunc("/api/holidays/{date}", handlers.DateHandler).Methods(http.MethodGet)

	return r
}
