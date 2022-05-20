package handler

import (
	"Douyin/global"
	"Douyin/proto"
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RelationAction(ctx *gin.Context) {
	token := ctx.Query("token")
	user_id := ctx.Query("user_id")
	to_user_id := ctx.Query("to_user_id")
	action_type := ctx.Query("action_type")

	uid, _ := strconv.Atoi(user_id)
	tuid, _ := strconv.Atoi(to_user_id)
	ac, _ := strconv.Atoi(action_type)

	RelationSrv := global.ConnMap[global.ServerConfig.SrvServerInfo.RelationSrv]
	rsp, err := RelationSrv.RelationAction(context.Background(), &proto.DouyinRelationActionRequest{
		UserId:     int64(uid),
		Token:      token,
		ToUserId:   int64(tuid),
		ActionType: int32(ac),
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

func FollowList(ctx *gin.Context) {
	token := ctx.Query("token")
	user_id := ctx.Query("user_id")

	uid, _ := strconv.Atoi(user_id)

	RelationSrv := global.ConnMap[global.ServerConfig.SrvServerInfo.RelationSrv]
	rsp, err := RelationSrv.RelationFollowList(context.Background(), &proto.DouyinRelationFollowListRequest{
		UserId: int64(uid),
		Token:  token,
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

func FollowerList(ctx *gin.Context) {
	token := ctx.Query("token")
	user_id := ctx.Query("user_id")

	uid, _ := strconv.Atoi(user_id)

	RelationSrv := global.ConnMap[global.ServerConfig.SrvServerInfo.RelationSrv]
	rsp, err := RelationSrv.RelationFollowerList(context.Background(), &proto.DouyinRelationFollowerListRequest{
		UserId: int64(uid),
		Token:  token,
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
