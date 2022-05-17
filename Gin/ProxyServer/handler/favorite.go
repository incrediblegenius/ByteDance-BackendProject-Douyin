package handler

import (
	"Douyin/ProxyServer/client"
	"Douyin/proto"
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func FavoriteAction(ctx *gin.Context) {
	token := ctx.Query("token")
	user_id := ctx.Query("user_id")
	video_id := ctx.Query("video_id")
	action_type := ctx.Query("action_type")
	uid, _ := strconv.Atoi(user_id)
	vid, _ := strconv.Atoi(video_id)
	ac, _ := strconv.Atoi(action_type)

	rsp, err := client.SrvClient.FavoriteAction(context.Background(), &proto.DouyinFavoriteActionRequest{
		Token:   token,
		UserId:  int64(uid),
		VideoId: int64(vid),
		Action:  int32(ac),
	})
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status_code": -1,
			"status_msg":  err.Error(),
		})
		return
	} else if rsp.StatusCode != 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status_code": rsp.StatusCode,
			"status_msg":  rsp.StatusMsg,
		})
		return
	}
	ctx.JSON(http.StatusOK, rsp)
}

func FavoriteList(ctx *gin.Context) {
	token := ctx.Query("token")
	user_id := ctx.Query("user_id")
	uid, _ := strconv.Atoi(user_id)

	rsp, err := client.SrvClient.FavoriteList(context.Background(), &proto.DouyinFavoriteListRequest{
		Token:  token,
		UserId: int64(uid),
	})
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status_code": -1,
			"status_msg":  err.Error(),
		})
		return
	} else if rsp.StatusCode != 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status_code": rsp.StatusCode,
			"status_msg":  rsp.StatusMsg,
		})
		return
	}
	ctx.JSON(http.StatusOK, rsp)
}
