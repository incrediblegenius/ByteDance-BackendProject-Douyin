package main

import (
	"Douyin/global"
	"os"
	"os/exec"
	"strings"
)

func main() {
	filename := "112121"
	cmd := []string{
		"$(docker run --rm -i -v",
		global.ServerConfig.StaticInfo.StaticDir + "/tmp:/tmp",
		"linuxserver/ffmpeg",
		"-i /tmp/test.mp4",
		"-ss 00:00:05",
		"-frames:v 1 test.png",
		"-c:a copy /tmp/test.png)",
	}
	exec.Command("/bin/bash", "-c", strings.Join(cmd, " ")).Run()

	_ = os.Rename(global.ServerConfig.StaticInfo.StaticDir+"/tmp/test.mp4", global.ServerConfig.StaticInfo.StaticDir+"/videos/"+filename+".mp4")
	_ = os.Rename(global.ServerConfig.StaticInfo.StaticDir+"/tmp/test.png", global.ServerConfig.StaticInfo.StaticDir+"/covers/"+filename+".png")

}
