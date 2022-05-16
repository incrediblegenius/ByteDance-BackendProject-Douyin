package main

import (
	"Douyin/cfg"
	"Douyin/model"
	"Douyin/user_srv/global"
	"bufio"
	"fmt"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	f, _ := os.Open("./data.txt")
	defer f.Close()
	input := bufio.NewScanner(f)
	cnt := 0
	rand.Seed(time.Now().UnixMilli())
	for input.Scan() {
		if (cnt+1)%100 == 0 {
			url := input.Text()
			// fmt.Println(url)
			ID := rand.Intn(1000) + 1
			result := global.DB.Create(&model.Video{
				AuthorID:   ID,
				PlayUrl:    url,
				CoverUrl:   url,
				IsFavorite: false,
			})
			if result.Error != nil {
				fmt.Println("插入失败")
			}
		}
		cnt++
	}
	// vs := []model.Video{}
	// global.DB.Find(&vs)
	// fmt.Println(len(vs))
}

func SaveVideoAndCover(data *multipart.FileHeader) error {
	var err error
	f, err := data.Open()
	if err != nil {
		return err
	}
	defer f.Close()
	buf := make([]byte, data.Size)
	_, err = f.Read(buf)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(cfg.StaticDir+"/tmp/test.mp4", buf, 0644)
	if err != nil {
		return err
	}
	cmd := []string{
		"$(docker run --rm -i -v",
		cfg.StaticDir + "/tmp:/tmp",
		"linuxserver/ffmpeg",
		"-i /tmp/test.mp4",
		"-ss 00:00:05",
		"-frames:v 1 test.png",
		"-c:a copy /tmp/test.png)",
	}
	err = exec.Command("/bin/bash", "-c", strings.Join(cmd, " ")).Run()
	if err != nil {
		return err
	}

	return nil
}
