package main

import (
	grpcserver "Douyin/user_srv/grpcServer"
	"flag"

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
	if err := grpcserver.Run(); err != nil {
		glog.Fatal(err)
	}
}
