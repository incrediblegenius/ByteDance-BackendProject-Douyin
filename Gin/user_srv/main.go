package main

import (
	"Douyin/cfg"
	grpcserver "Douyin/user_srv/grpcServer"
	"flag"
	"net/http"
	"os"

	"github.com/golang/glog"
)

func main() {
	flag.Parse()
	defer glog.Flush()
	//go func() {
	//	err := grpcserver.Run()
	//	if err != nil {
	//		glog.Fatal(err)
	//	}
	//}()
	go func() {
		os.MkdirAll(cfg.StaticDir+"/videos", 0777)
		os.MkdirAll(cfg.StaticDir+"/covers", 0777)
		os.MkdirAll(cfg.StaticDir+"/tmp", 0777)
		http.Handle("/", http.FileServer(http.Dir(cfg.StaticDir)))
		if e := http.ListenAndServe(":8081", nil); e != nil {
			glog.Fatal(e)
		}
	}()
	go func() {
		if err := grpcserver.VideoServerRun(); err != nil {
			glog.Fatal(err)
		}
	}()
	if err := grpcserver.Run(); err != nil {
		glog.Fatal(err)
	}
}