package main

import (
	"UserServer/cfg"
	"UserServer/model"
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

const (
	Dir = "/Users/evil/Desktop/Go/Douyin/Gin/model/FakeData"
)

var (
	DB *gorm.DB
)

func init() {
	addr := "root:root@tcp(localhost:3306)/douyin_user?charset=utf8mb4&parseTime=True&loc=Local"
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	DB, _ = gorm.Open(mysql.Open(addr), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})
	rand.Seed(time.Now().UnixMilli())
}

func CreateVideos() {
	f, _ := os.Open("./data.txt")
	defer f.Close()
	input := bufio.NewScanner(f)
	cnt := 0

	ch := make(chan struct{}, 10)
	wg := &sync.WaitGroup{}
	for input.Scan() {
		if (cnt+2)%200 == 0 {
			// 并发下载，之后写（需要考虑文件io的并发）
			url := input.Text()
			ch <- struct{}{}
			wg.Add(1)
			go func(url string, cnt int) {
				ID := rand.Intn(50) + 1

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
				result := DB.Create(&model.Video{
					AuthorID: ID,
					PlayUrl:  url,
					CoverUrl: fmt.Sprintf("http://%s:%d/covers/%s.png", cfg.ServerIP, cfg.ServerPort, filename),
				})
				if result.Error != nil {
					fmt.Println("插入失败")
				}
				os.Rename(fmt.Sprintf(Dir+"/test%d.png", cnt), cfg.StaticDir+"/covers/"+filename+".png")
				os.Remove(fmt.Sprintf(Dir+"/test%d.mp4", cnt))
				<-ch
				wg.Done()
			}(url, cnt)
		}
		cnt++
	}
	wg.Wait()
	// vs := []model.Video{}
	// global.DB.Find(&vs)
	// fmt.Println(len(vs))
}

func SaveVideoAndCover(cnt int) error {
	cmd := []string{
		"$(docker run --rm -i -v",
		Dir + ":/tmp",
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

func CreateLikes(nums int) {

	for i := 0; i < nums; i++ {
		uid := rand.Intn(50) + 1
		vid := rand.Intn(228) + 1

		like := model.FavoriteVideo{}
		result := DB.First(&like, "user_id = ? and video_id = ?", uid, vid)
		if result.RowsAffected == 0 {
			DB.Transaction(func(tx *gorm.DB) error {
				// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
				if err := tx.Create(&model.FavoriteVideo{UserID: uid, VideoID: vid}).Error; err != nil {
					// 返回任何错误都会回滚事务
					return err
				}
				if err := tx.Model(&model.Video{}).Where("id = ?", vid).Update("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error; err != nil {
					return err
				}
				// 返回 nil 提交事务
				return nil
			})
		}
	}
}

func CreateRealations(nums int) {
	for i := 0; i < nums; i++ {
		uid := rand.Intn(50) + 1
		tuid := rand.Intn(50) + 1
		DB.Transaction(func(tx *gorm.DB) error {
			// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
			if err := tx.Create(&model.Relation{
				FollowFrom: int(uid),
				FollowTo:   int(tuid),
			}).Error; err != nil {
				// 返回任何错误都会回滚事务
				return err
			}
			if err := tx.Model(&model.User{}).Where("id = ?", uid).Update("following_count", gorm.Expr("following_count + ?", 1)).Error; err != nil {
				return err
			}
			if err := tx.Model(&model.User{}).Where("id = ?", tuid).Update("follower_count", gorm.Expr("follower_count + ?", 1)).Error; err != nil {
				return err
			}
			// 返回 nil 提交事务
			return nil
		})
	}
}

func CountComments() {
	videos := []model.Video{}
	DB.Find(&videos)
	for _, v := range videos {
		var cnt int64
		DB.Model(&model.Comment{}).Where("video_id = ?", v.ID).Count(&cnt)
		DB.Model(&v).Update("comment_count", cnt)
	}
}

func main() {
	// CreateVideos()
	// CreateLikes(1000)
	// CreateRealations(1000)
	CountComments()
}
