package main

import (
	"github.com/xiaoxuan6/chinese-holidays-api/router"
	"log"
	"net/http"
)

func main() {
	r := router.InitRouter()
	log.Panic(http.ListenAndServe(":80", r))
}
