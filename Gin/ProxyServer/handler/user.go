package handler

import (
	"Douyin/ProxyServer/client"
	"Douyin/proto/userproto"
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Register(ctx *gin.Context) {
	username := ctx.Query("username")
	password := ctx.Query("password")
	if username == "" || password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status_code": 1,
			"status_msg":  "username or password is empty",
		})
		return
	}
	var rsp *userproto.DouyinUserRegisterResponse
	rsp, _ = client.UserSrvClient.Register(context.Background(), &userproto.DouyinUserRegisterRequest{
		Username: username,
		Password: password,
	})
	if rsp.StatusCode != 0 {
		//fmt.Println(rsp)
		ctx.JSON(http.StatusOK, gin.H{
			"status_code": rsp.StatusCode,
			"status_msg":  rsp.StatusMsg,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status_code": rsp.StatusCode,
		"status_msg":  rsp.StatusMsg,
		"user_id":     rsp.UserId,
		"token":       rsp.Token,
	})
}

func Login(ctx *gin.Context) {
	username := ctx.Query("username")
	password := ctx.Query("password")
	if username == "" || password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status_code": -1,
			"status_msg":  "username or password is empty",
		})
		return
	}
	var rsp *userproto.DouyinUserRegisterResponse
	rsp, _ = client.UserSrvClient.Login(context.Background(), &userproto.DouyinUserRegisterRequest{
		Username: username,
		Password: password,
	})
	if rsp.StatusCode != 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"status_code": rsp.StatusCode,
			"status_msg":  rsp.StatusMsg,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status_code": rsp.StatusCode,
		"status_msg":  rsp.StatusMsg,
		"user_id":     rsp.UserId,
		"token":       rsp.Token,
	})
}

func GetUserInfo(ctx *gin.Context) {
	userId := ctx.Query("user_id")
	token := ctx.Query("token")
	//fmt.Println(userId, token)
	id, err := strconv.Atoi(userId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status_code": -1,
			"status_msg":  err.Error(),
		})
		return
	}
	rsp, err := client.UserSrvClient.GetUserById(context.Background(), &userproto.IdRequest{
		Id:        int64(id),
		Token:     token,
		NeedToken: true,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status_code": -1,
			"status_msg":  err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status_code": 0,
		"status_msg":  "success",
		"user":        rsp,
	})
}
