package handler

import (
	"CommentSrv/global"
	"CommentSrv/model"
	"CommentSrv/proto"
	"context"
	"errors"

	"gorm.io/gorm"
)

type Server struct {
	proto.UnimplementedServerServer
}

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
			if err := tx.Unscoped().Delete(&model.Comment{}, req.CommentId).Error; err != nil {
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
	cid := claim.Id
	if uid == 0 {
		uid = cid
	}
	var comments []*model.Comment
	global.DB.Where("video_id = ?", req.VideoId).Find(&comments)
	commentsList := make([]*proto.Comment, len(comments))
	for i := range comments {
		user, _ := GetUserById(&Request{
			Id:       int64(comments[i].UserID),
			SearchId: cid,
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

type Request struct {
	Id       int64
	SearchId int64
}

func GetUserById(req *Request) (*proto.User, error) {

	var user model.User
	ans := &proto.User{}
	result := global.DB.First(&user, req.Id)
	if result.RowsAffected == 0 {
		return nil, errors.New("未找到用户")
	}
	ans.Name = user.UserName
	ans.Id = int64(user.ID)
	ans.FollowCount = int64(user.FollowingCount)
	ans.FollowerCount = int64(user.FollowerCount)
	if req.SearchId == 0 {
		ans.IsFollow = false
	} else {
		result := global.DB.Where("follow_from = ? AND follow_to = ?", req.SearchId, req.Id).First(&model.Relation{})
		if result.RowsAffected != 0 {
			ans.IsFollow = true
		} else {
			ans.IsFollow = false
		}
	}

	return ans, nil
}
