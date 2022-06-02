package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/namsral/flag"
)

var (
	url        string
	concurency int
	rounds     int
	sleep      int
)

func init() {
	flag.IntVar(&concurency, "c", 1, "concurency(max running routines)")
	flag.IntVar(&rounds, "r", 1, "rounds(total requests)")
	flag.IntVar(&sleep, "s", 0, "sleep(ms)")
	flag.StringVar(&url, "url", "https://www.baidu.com", "test url")
	flag.Parse()
}

func main() {
	fmt.Println("start")
	fmt.Println("url:", url)
	fmt.Println("concurency:", concurency)
	fmt.Println("rounds:", rounds)
	fmt.Println("sleep:", sleep)
	wg := sync.WaitGroup{}
	syncChan := make(chan struct{}, concurency)
	start := time.Now().UnixNano()
	cnt := 0 // 计数器
	for i := 0; i < rounds; i++ {
		if sleep != 0 {
			time.Sleep(time.Millisecond * time.Duration(sleep))
		}
		wg.Add(1)
		syncChan <- struct{}{}
		cnt += 1
		go func(cnt int) {
			defer wg.Done()
			defer func() {
				<-syncChan
			}()
			resp, err := http.Get(url)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("No.", cnt, "StatusCode:", resp.StatusCode)
			defer resp.Body.Close()
			_, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
				return
			}
		}(cnt)
	}
	wg.Wait()
	end := time.Now().UnixNano()
	fmt.Println("time cost(nano)", end-start)
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
