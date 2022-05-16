package handler

import (
	"Douyin/ProxyServer/userClient"
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
	rsp, _ := userClient.UserSrvClient.GetUserFeed(context.Background(), &proto.DouyinFeedRequest{
		LatestTime: t,
		Token:      token,
	})
	// fmt.Println(rsp)

	ctx.JSON(http.StatusOK, rsp)

}
