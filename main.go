package main

import (
    "github.com/kelseyhightower/envconfig"
    "github.com/xiaoxuan6/chinese-holidays-api/router"
    "log"
    "net/http"
    "strconv"
)

type env struct {
    PORT int
}

var e env

func init() {
    if err := envconfig.Process("holiday", &e); err != nil {
        log.Println("envconfig load fail:", err.Error())
    }

    if e.PORT == 0 {
        e.PORT = 80
    }
}

func main() {
    r := router.InitRouter()

    port := strconv.Itoa(e.PORT)
    log.Println("已启动服务：127.0.0.1:" + port)
    if err := http.ListenAndServe(":"+port, r); err != nil {
        log.Panic(err)
    }
}
