package api

import (
	"github.com/xiaoxuan6/chinese-holidays-api/router"
	"net/http"
)

var routers http.Handler

func init() {
	routers = router.InitRouter()
}

func Api(w http.ResponseWriter, r *http.Request) {
	routers.ServeHTTP(w, r)
}
