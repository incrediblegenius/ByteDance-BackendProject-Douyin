package main

import (
	"Douyin/model"
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

const (
	Dir = "/Users/evil/Desktop/Go/Douyin/Gin/model/FakeData"
)

var (
	DB        *gorm.DB
	OssClient *cos.Client
)

func init() {
	addr := "root:root@tcp(113.54.157.200:3307)/douyin_user?charset=utf8mb4&parseTime=True&loc=Local"
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
	u, _ := url.Parse("https://doiuyin-1302721364.cos.ap-chengdu.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	OssClient = cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			// 通过环境变量获取密钥
			// 环境变量 SECRETID 表示用户的 SecretId，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretID: "AKIDAIC1cy62k7HDwQfhU4PWO32xhGgtvlOp",
			// 环境变量 SECRETKEY 表示用户的 SecretKey，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretKey: "GI7lCPRIxPfjcIl14vZ3MTN4ZqsgI0Xs",
		},
	})
}

func CreateVideos() {
	f, _ := os.Open("./data.txt")
	defer f.Close()
	input := bufio.NewScanner(f)
	cnt := 0

	ch := make(chan struct{}, 30)
	wg := &sync.WaitGroup{}
	for input.Scan() {
		if (cnt+1)%11 == 0 {
			// 并发下载，之后写（需要考虑文件io的并发）
			url := input.Text()
			ch <- struct{}{}
			wg.Add(1)
			go func(url string, cnt int) {
				ID := rand.Intn(50) + 1

				urlSlice := strings.Split(url, "/")
				tmp := urlSlice[len(urlSlice)-1]
				filename := tmp[:len(tmp)-4]
				out, err := os.Create(tmp)
				if err != nil {
					fmt.Println(err)
					return
				}
				defer out.Close()
				resp, err := http.Get(url)
				if err != nil || resp.StatusCode != 200 {
					fmt.Println(err)
					resp.Body.Close()
					return
				}
				defer resp.Body.Close()
				_, err = io.Copy(out, resp.Body)
				if err != nil {
					fmt.Println(err)
					return
				}
				// mutex.Lock()
				err = SaveVideoAndCover(filename)
				// mutex.Unlock()

				if err == nil {
					result := DB.Create(&model.Video{
						AuthorID: ID,
						PlayUrl:  url,
						CoverUrl: fmt.Sprintf("https://doiuyin-1302721364.cos.ap-chengdu.myqcloud.com/covers/%s.png", filename),
					})
					if result.Error != nil {
						fmt.Println("插入失败")
					}
				}
				<-ch
				defer wg.Done()
			}(url, cnt)
		}
		cnt++
	}
	wg.Wait()
	// vs := []model.Video{}
	// global.DB.Find(&vs)
	// fmt.Println(len(vs))
}

func SaveVideoAndCover(filename string) error {
	cmd := []string{
		"$(docker run --rm -i -v",
		Dir + ":/tmp",
		"linuxserver/ffmpeg",
		fmt.Sprintf("-i /tmp/%s.mp4", filename),
		"-ss 00:00:05",
		"-frames:v 1 -vf scale=iw/4:ih/4 ",
		fmt.Sprintf("/tmp/%s.png)", filename),
	}

	c := exec.Command("/bin/bash", "-c", strings.Join(cmd, " "))
	fmt.Println(c.String())
	err := c.Run()
	// fmt.Println("ffmpeg")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	vkey := "/covers/" + filename + ".png"
	_, _, err = OssClient.Object.Upload(
		context.Background(), vkey, fmt.Sprintf("%s/%s.png", Dir, filename), nil,
	)
	if err != nil {
		fmt.Println("upload video error:", err)
		return err
	}
	defer os.Remove(fmt.Sprintf("%s/%s.mp4", Dir, filename))
	defer os.Remove(fmt.Sprintf("%s/%s.png", Dir, filename))

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
func CountFavorite() {
	videos := []model.Video{}
	DB.Find(&videos)
	for _, v := range videos {
		var cnt int64
		DB.Model(&model.FavoriteVideo{}).Where("video_id = ?", v.ID).Count(&cnt)
		DB.Model(&v).Update("favorite_count", cnt)
	}
}

func main() {
	// CreateVideos()
	// CreateLikes(1000)
	// CreateRealations(1000)
	// CountComments()
	CountFavorite()
}
