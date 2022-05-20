package handler

import (
	"Douyin/global"
	"Douyin/proto"
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetFeed(ctx *gin.Context) {
	latest_time := ctx.Query("latest_time")
	token := ctx.Query("token")
	var t int64
	if latest_time != "" {
		tmp, _ := strconv.Atoi(latest_time)
		t = int64(tmp)
	} else {
		t = time.Now().UnixMilli()
	}
	// fmt.Println(t)
	FeedSrv := global.ConnMap[global.ServerConfig.SrvServerInfo.FeedSrv]
	rsp, err := FeedSrv.GetUserFeed(context.Background(), &proto.DouyinFeedRequest{
		LatestTime: t,
		Token:      token,
	})
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status_code": -1,
			"status_msg":  err.Error(),
		})
		return
	}
	// fmt.Println(rsp)

	ctx.JSON(http.StatusOK, rsp)

}
