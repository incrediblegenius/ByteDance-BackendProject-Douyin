package main

import (
	"Douyin/cfg"
	"os"
	"os/exec"
	"strings"
)

func main() {
	filename := "112121"
	cmd := []string{
		"$(docker run --rm -i -v",
		cfg.StaticDir + "/tmp:/tmp",
		"linuxserver/ffmpeg",
		"-i /tmp/test.mp4",
		"-ss 00:00:05",
		"-frames:v 1 test.png",
		"-c:a copy /tmp/test.png)",
	}
	exec.Command("/bin/bash", "-c", strings.Join(cmd, " ")).Run()

	_ = os.Rename(cfg.StaticDir+"/tmp/test.mp4", cfg.StaticDir+"/videos/"+filename+".mp4")
	_ = os.Rename(cfg.StaticDir+"/tmp/test.png", cfg.StaticDir+"/covers/"+filename+".png")

}
