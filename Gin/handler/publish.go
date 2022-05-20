package handler

import (
	"Douyin/global"
	"strconv"

	"Douyin/proto"
	"context"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Publish(ctx *gin.Context) {
	token := ctx.PostForm("token")
	data, err := ctx.FormFile("data")
	title := ctx.Query("title")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status_code": -1,
			"status_msg":  err.Error(),
		})
		return
	}
	filename := fmt.Sprintf("%s_%d", token[:10], time.Now().UnixMilli())
	c := make(chan error)
	go func(data *multipart.FileHeader) {
		err := SaveVideoAndCover(data)
		c <- err
	}(data)

	PublishSrv := global.ConnMap[global.ServerConfig.SrvServerInfo.PublishSrv]
	rsp, err := PublishSrv.PublishAction(context.Background(), &proto.DouyinPublishActionRequest{
		Token:     token,
		VideoName: filename,
		Title:     title,
	})

	if e := <-c; e == nil && err == nil && rsp.StatusCode == 0 {
		os.Rename(global.ServerConfig.StaticInfo.StaticDir+"/tmp/test.mp4", global.ServerConfig.StaticInfo.StaticDir+"/videos/"+filename+".mp4")
		os.Rename(global.ServerConfig.StaticInfo.StaticDir+"/tmp/test.png", global.ServerConfig.StaticInfo.StaticDir+"/covers/"+filename+".png")
		ctx.JSON(http.StatusOK, rsp)
	} else {
		os.Remove(global.ServerConfig.StaticInfo.StaticDir + "/tmp/test.mp4")
		os.Remove(global.ServerConfig.StaticInfo.StaticDir + "/tmp/test.png")
		if e != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status_code": -1,
				"status_msg":  e.Error(),
			})
		} else {
			ctx.JSON(http.StatusBadRequest, rsp)
		}
	}
}

func PublishList(ctx *gin.Context) {
	token := ctx.Query("token")
	user_id := ctx.Query("user_id")
	uid, _ := strconv.Atoi(user_id)
	PublishSrv := global.ConnMap[global.ServerConfig.SrvServerInfo.PublishSrv]
	rsp, err := PublishSrv.PublishList(context.Background(), &proto.DouyinPublishListRequest{
		Token:  token,
		UserId: int64(uid),
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}
	// fmt.Println(rsp)
	ctx.JSON(http.StatusOK, rsp)
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
	err = ioutil.WriteFile(global.ServerConfig.StaticInfo.StaticDir+"/tmp/test.mp4", buf, 0644)
	if err != nil {
		return err
	}
	cmd := []string{
		"$(docker run --rm -i -v",
		global.ServerConfig.StaticInfo.StaticDir + "/tmp:/tmp",
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
