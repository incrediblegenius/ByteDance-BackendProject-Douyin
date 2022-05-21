package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

var (
	Geturl = []string{
		"http://127.0.0.1:8080/douyin/user/?user_id=51&token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6NTEsIkF1dGhvcml0eUlkIjowfQ.u8q99T62bLy8fH-TvtT1C--aP5OKBb3h_8UuJyoGIaU",
		"http://127.0.0.1:8080/douyin/relation/follow/list/?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6NTEsIkF1dGhvcml0eUlkIjowfQ.u8q99T62bLy8fH-TvtT1C--aP5OKBb3h_8UuJyoGIaU&user_id=51",
		"http://127.0.0.1:8080/douyin/publish/list/?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6NTEsIkF1dGhvcml0eUlkIjowfQ.u8q99T62bLy8fH-TvtT1C--aP5OKBb3h_8UuJyoGIaU&user_id=47",
		"http://10.252.138.8:8080/douyin/feed?latest_time=1652864872162&token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6NTEsIkF1dGhvcml0eUlkIjowfQ.u8q99T62bLy8fH-TvtT1C--aP5OKBb3h_8UuJyoGIaU",
		"http://10.252.138.8:8080/douyin/favorite/list/?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6NTEsIkF1dGhvcml0eUlkIjowfQ.u8q99T62bLy8fH-TvtT1C--aP5OKBb3h_8UuJyoGIaU&user_id=51",
		"http://127.0.0.1:8080/douyin/comment/list/?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6NTEsIkF1dGhvcml0eUlkIjowfQ.u8q99T62bLy8fH-TvtT1C--aP5OKBb3h_8UuJyoGIaU&video_id=118",
	}
)

func main() {
	for _, url := range Geturl {
		wg := sync.WaitGroup{}
		syncChan := make(chan struct{}, 15000)
		start := time.Now().Unix()
		for i := 0; i < 20000; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				resp, err := http.Get(url)
				if err != nil {
					fmt.Println(err)
					return
				}
				defer resp.Body.Close()
				_, err = ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Println(err)
					return
				}
				syncChan <- struct{}{}
			}()
		}
		wg.Wait()
		end := time.Now().Unix()
		fmt.Println(len(syncChan), "time cost", end-start)
	}
	// resp, err := http.Post(url, "application/json", nil)
	// if err != nil {
	// 	panic(err)
	// }

	// defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(len(body))
}
