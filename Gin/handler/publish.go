package handler

import (
	"Douyin/global"
	"strconv"

	"Douyin/proto"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Publish(ctx *gin.Context) {
	token := ctx.PostForm("token")
	data, err := ctx.FormFile("data")
	title := ctx.Query("title")
	f, err := data.Open()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status_code": -1,
			"status_msg":  err.Error(),
		})
		return
	}
	defer f.Close()
	buf := make([]byte, data.Size)
	_, err = f.Read(buf)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status_code": -1,
			"status_msg":  err.Error(),
		})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status_code": -1,
			"status_msg":  err.Error(),
		})
		return
	}
	filename := fmt.Sprintf("%s_%d", token[:10], time.Now().UnixMilli())

	PublishSrv := global.ConnMap["publish_srv"]
	rsp, err := PublishSrv.PublishAction(context.Background(), &proto.DouyinPublishActionRequest{
		Token:     token,
		Content:   buf,
		VideoName: filename,
		Title:     title,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}
	ctx.JSON(http.StatusOK, rsp)
}

func PublishList(ctx *gin.Context) {
	token := ctx.Query("token")
	user_id := ctx.Query("user_id")
	uid, _ := strconv.Atoi(user_id)
	PublishSrv := global.ConnMap["publish_srv"]
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
