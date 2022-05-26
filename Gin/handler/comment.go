package handler

import (
	"Douyin/global"
	"Douyin/proto"
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CommentAction(ctx *gin.Context) {
	token := ctx.Query("token")
	user_id := ctx.Query("user_id")
	video_id := ctx.Query("video_id")
	action_type := ctx.Query("action_type")
	comment_text := ctx.Query("comment_text")
	comment_id := ctx.Query("comment_id")
	uid, _ := strconv.Atoi(user_id)
	vid, _ := strconv.Atoi(video_id)
	ac, _ := strconv.Atoi(action_type)
	cid := 0
	if ac == 2 {
		cid, _ = strconv.Atoi(comment_id)
	}
	userSrv := global.ConnMap["user_srv"]
	rsp, err := userSrv.CommentAction(context.Background(), &proto.DouyinCommentActionRequest{
		UserId:      int64(uid),
		Token:       token,
		VideoId:     int64(vid),
		ActionType:  int32(ac),
		CommentText: comment_text,
		CommentId:   int64(cid),
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

func CommentList(ctx *gin.Context) {
	token := ctx.Query("token")
	user_id := ctx.Query("user_id")
	video_id := ctx.Query("video_id")
	vid, _ := strconv.Atoi(video_id)
	uid, _ := strconv.Atoi(user_id)
	CommentSrv := global.ConnMap["comment_srv"]
	rsp, err := CommentSrv.CommentList(context.Background(), &proto.DouyinCommentListRequest{
		UserId:  int64(uid),
		Token:   token,
		VideoId: int64(vid),
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
