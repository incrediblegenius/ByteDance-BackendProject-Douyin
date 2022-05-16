package main

import (
	"Douyin/cfg"
	"Douyin/model"
	"Douyin/user_srv/global"
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"net/http"
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
	ch := make(chan struct{}, 10)
	// mutex := &sync.Mutex{}
	for input.Scan() {
		if (cnt+1)%100 == 0 {
			// 并发下载，之后写（需要考虑文件io的并发）
			url := input.Text()
			ch <- struct{}{}
			go func(url string, cnt int) {
				ID := rand.Intn(1000) + 1

				urlSlice := strings.Split(url, "/")
				tmp := urlSlice[len(urlSlice)-1]
				filename := tmp[:len(tmp)-4]
				out, err := os.Create(fmt.Sprintf("./test%d.mp4", cnt))
				if err != nil {
					fmt.Println(err)
					return
				}
				defer out.Close()
				resp, err := http.Get(url)
				if err != nil {
					fmt.Println(err)
					return
				}
				_, err = io.Copy(out, resp.Body)
				if err != nil {
					fmt.Println(err)
					return
				}
				// mutex.Lock()
				SaveVideoAndCover(cnt)
				// mutex.Unlock()
				defer resp.Body.Close()
				result := global.DB.Create(&model.Video{
					AuthorID:   ID,
					PlayUrl:    url,
					CoverUrl:   fmt.Sprintf("http://%s:%d/covers/%s.png", cfg.ServerIP, cfg.ServerPort, filename),
					IsFavorite: false,
				})
				if result.Error != nil {
					fmt.Println("插入失败")
				}
				os.Rename(fmt.Sprintf("/Users/evil/Desktop/Go/Douyin/model/FakeData/test%d.png", cnt), cfg.StaticDir+"/covers/"+filename+".png")
				os.Remove(fmt.Sprintf("/Users/evil/Desktop/Go/Douyin/model/FakeData/test%d.mp4", cnt))
				<-ch
			}(url, cnt)
		}
		cnt++
	}
	// vs := []model.Video{}
	// global.DB.Find(&vs)
	// fmt.Println(len(vs))
}

func SaveVideoAndCover(cnt int) error {
	cmd := []string{
		"$(docker run --rm -i -v",
		"/Users/evil/Desktop/Go/Douyin/model/FakeData:/tmp",
		"linuxserver/ffmpeg",
		fmt.Sprintf("-i /tmp/test%d.mp4", cnt),
		"-ss 00:00:05",
		"-frames:v 1 test.png",
		fmt.Sprintf("-c:a copy /tmp/test%d.png)", cnt),
	}
	err := exec.Command("/bin/bash", "-c", strings.Join(cmd, " ")).Run()
	if err != nil {
		return err
	}
	return nil
}
