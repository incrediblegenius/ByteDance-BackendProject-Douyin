package main

import (
	_ "Douyin/ProxyServer/client"
	"Douyin/ProxyServer/router"
	"Douyin/cfg"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

func main() {
	r := gin.Default()
	go func() {
		os.MkdirAll(cfg.StaticDir+"/videos", 0777)
		os.MkdirAll(cfg.StaticDir+"/covers", 0777)
		os.MkdirAll(cfg.StaticDir+"/tmp", 0777)
		http.Handle("/", http.FileServer(http.Dir(cfg.StaticDir)))
		if e := http.ListenAndServe(":8081", nil); e != nil {
			glog.Fatal(e)
		}
	}()
	router.InitRouter(r)
	err := r.Run(":8080")
	if err != nil {
		return
	}
}
