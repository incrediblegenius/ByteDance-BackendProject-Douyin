package handler

import (
	"Douyin/model"
	"Douyin/proto/videoproto"
	"Douyin/user_srv/global"
	"context"
	"errors"

	"gorm.io/gorm"
)

type VideosServer struct {
	videoproto.UnimplementedVideosServer
}

func (s *VideosServer) FavoriteAction(ctx context.Context, req *videoproto.DouyinFavoriteActionRequest) (*videoproto.DouyinFavoriteActionResponse, error) {
	claim, err := global.Jwt.ParseToken(req.Token)
	if err != nil {
		return &videoproto.DouyinFavoriteActionResponse{
			StatusCode: -2,
			StatusMsg:  "token鉴权失败",
		}, nil
	}
	uid := req.UserId
	if uid == 0 {
		uid = claim.Id
	}
	like := model.FavoriteVideo{}
	result := global.DB.First(&like, "user_id = ? and video_id = ?", uid, req.VideoId)
	if result.RowsAffected == 0 && req.Action == 1 {
		global.DB.Transaction(func(tx *gorm.DB) error {
			// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
			if err := tx.Create(&model.FavoriteVideo{UserID: int(uid), VideoID: int(req.VideoId)}).Error; err != nil {
				// 返回任何错误都会回滚事务
				return err
			}
			if err := tx.Model(&model.Video{}).Where("id = ?", req.VideoId).Update("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error; err != nil {
				return err
			}
			// 返回 nil 提交事务
			return nil
		})
	} else if result.RowsAffected == 1 && req.Action == 2 {
		global.DB.Transaction(func(tx *gorm.DB) error {
			// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
			if err := tx.Delete(&like).Error; err != nil {
				// 返回任何错误都会回滚事务
				return err
			}
			if err := tx.Model(&model.Video{}).Where("id = ?", req.VideoId).Update("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error; err != nil {
				return err
			}
			// 返回 nil 提交事务
			return nil
		})
	} else {
		return &videoproto.DouyinFavoriteActionResponse{
			StatusCode: -1,
			StatusMsg:  "操作请求与实际不符",
		}, nil
	}
	return &videoproto.DouyinFavoriteActionResponse{
		StatusCode: 0,
		StatusMsg:  "操作成功",
	}, nil
}

func (s *VideosServer) GetUserById(ctx context.Context, req *videoproto.OthersRequest) (*videoproto.User, error) {
	ans := &videoproto.User{}
	var user model.User
	result := global.DB.First(&user, req.Id)
	if result.RowsAffected == 0 {
		return nil, errors.New("未找到用户")
	}
	ans.Name = user.UserName
	ans.Id = int64(user.ID)
	var cnt int64
	global.DB.Model(&model.Relation{}).Where("follow_from = ?", req.Id).Count(&cnt)
	ans.FollowCount = cnt
	global.DB.Model(&model.Relation{}).Where("follow_to = ?", req.Id).Count(&cnt)
	ans.FollowerCount = cnt
	if result := global.DB.First(&model.Relation{}, "follow_from = ? and follow_to = ?", req.CheckFrom, req.Id); result.RowsAffected == 1 {
		ans.IsFollow = true
	} else {
		ans.IsFollow = false
	}
	return ans, nil
}

// func (s *VideosServer) FavoriteList(ctx context.Context, req *videoproto.DouyinFavoriteListRequest) (*videoproto.DouyinFavoriteListResponse, error) {
// 	claim, err := global.Jwt.ParseToken(req.Token)
// 	if err != nil {
// 		return &videoproto.DouyinFavoriteListResponse{
// 			StatusCode: -2,
// 			StatusMsg:  "token鉴权失败",
// 		}, nil
// 	}
// 	uid := req.UserId
// 	if uid == 0 {
// 		uid = claim.Id
// 	}

// }
