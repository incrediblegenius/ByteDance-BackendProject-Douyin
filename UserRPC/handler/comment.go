package handler

import (
	"UserServer/global"
	"UserServer/model"
	"UserServer/proto"
	"context"

	"gorm.io/gorm"
)

func (s *Server) CommentAction(ctx context.Context, req *proto.DouyinCommentActionRequest) (*proto.DouyinCommentActionResponse, error) {
	claim, err := global.Jwt.ParseToken(req.Token)
	if err != nil {
		return &proto.DouyinCommentActionResponse{
			StatusCode: -2,
			StatusMsg:  "token鉴权失败",
		}, nil
	}
	uid := req.UserId
	if uid == 0 {
		uid = claim.Id
	}
	if req.ActionType == 1 && uid != 0 && req.VideoId != 0 {
		result := global.DB.Transaction(func(tx *gorm.DB) error {
			// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
			if err := tx.Create(&model.Comment{
				UserID:  int(uid),
				VideoID: int(req.VideoId),
				Content: req.CommentText,
			}).Error; err != nil {
				// 返回任何错误都会回滚事务
				return err
			}
			if err := tx.Model(&model.Video{}).Where("id = ?", req.VideoId).Update("comment_count", gorm.Expr("comment_count + ?", 1)).Error; err != nil {
				return err
			}
			// 返回 nil 提交事务
			return nil
		})
		if result != nil {
			return &proto.DouyinCommentActionResponse{
				StatusCode: -1,
				StatusMsg:  "操作失败",
			}, nil
		} else {
			return &proto.DouyinCommentActionResponse{
				StatusCode: 0,
				StatusMsg:  "发布评论成功",
			}, nil
		}
	} else if req.ActionType == 2 && uid != 0 && req.VideoId != 0 && req.CommentId != 0 {
		result := global.DB.Transaction(func(tx *gorm.DB) error {
			// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
			if err := tx.Delete(&model.Comment{}, req.CommentId).Error; err != nil {
				// 返回任何错误都会回滚事务
				return err
			}
			if err := tx.Model(&model.Video{}).Where("id = ?", req.VideoId).Update("comment_count", gorm.Expr("comment_count - ?", 1)).Error; err != nil {
				return err
			}
			// 返回 nil 提交事务
			return nil
		})
		if result != nil {
			return &proto.DouyinCommentActionResponse{
				StatusCode: -1,
				StatusMsg:  "操作失败",
			}, nil
		} else {
			return &proto.DouyinCommentActionResponse{
				StatusCode: 0,
				StatusMsg:  "删除评论成功",
			}, nil
		}
	}
	return &proto.DouyinCommentActionResponse{
		StatusCode: -3,
		StatusMsg:  "参数错误",
	}, nil
}
func (s *Server) CommentList(ctx context.Context, req *proto.DouyinCommentListRequest) (*proto.DouyinCommentListResponse, error) {
	claim, err := global.Jwt.ParseToken(req.Token)
	if err != nil {
		return &proto.DouyinCommentListResponse{
			StatusCode:  -2,
			StatusMsg:   "token鉴权失败",
			CommentList: []*proto.Comment{},
		}, nil
	}
	uid := req.UserId
	if uid == 0 {
		uid = claim.Id
	}
	var comments []*model.Comment
	global.DB.Where("user_id = ? and video_id = ?", uid, req.VideoId).Find(&comments)
	commentsList := make([]*proto.Comment, len(comments))
	for i := range comments {
		user, _ := s.GetUserById(context.Background(), &proto.IdRequest{
			Id:        int64(comments[i].UserID),
			NeedToken: false,
			SearchId:  uid,
		})
		commentsList[i] = &proto.Comment{
			Id:         int64(comments[i].ID),
			User:       user,
			Content:    comments[i].Content,
			CreateDate: comments[i].CreatedAt.Format("01-02"),
		}
	}
	return &proto.DouyinCommentListResponse{
		StatusCode:  0,
		StatusMsg:   "获取评论列表成功",
		CommentList: commentsList,
	}, nil
}
