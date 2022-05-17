package main

import (
	"UserServer/server"
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

	if err := server.Run(); err != nil {
		glog.Fatal(err)
	}
}
