package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

const (
	url = "http://127.0.0.1:8080/douyin/user/login/?username=zhou&password=123456"
)

func main() {
	wg := sync.WaitGroup{}
	syncChan := make(chan struct{}, 10000)
	start := time.Now().Unix()
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := http.Post(url, "application/json", nil)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
				return
			}
			if len(body) == 179 {
				syncChan <- struct{}{}
			}
		}()
		time.Sleep(time.Millisecond)
	}
	wg.Wait()
	end := time.Now().Unix()
	fmt.Println(len(syncChan), "time cost", end-start)

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
