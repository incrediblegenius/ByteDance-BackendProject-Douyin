package main

import (
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
		os.MkdirAll("./videos", 0777)
		http.Handle("/", http.FileServer(http.Dir("./videos")))
		if e := http.ListenAndServe(":8081", nil); e != nil {
			glog.Fatal(e)
		}
	}()
	if err := grpcserver.Run(); err != nil {
		glog.Fatal(err)
	}
}
