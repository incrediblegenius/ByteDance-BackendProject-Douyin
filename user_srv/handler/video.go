package handler

import (
	"Douyin/model"
	"Douyin/proto"
	"Douyin/user_srv/global"
	"context"

	"gorm.io/gorm"
)

type VideosServer struct {
	proto.UnimplementedVideosServer
}

func (s *VideosServer) FavoriteAction(ctx context.Context, req *proto.DouyinFavoriteActionRequest) (*proto.DouyinFavoriteActionResponse, error) {
	claim, err := global.Jwt.ParseToken(req.Token)
	if err != nil {
		return &proto.DouyinFavoriteActionResponse{
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
		return &proto.DouyinFavoriteActionResponse{
			StatusCode: -1,
			StatusMsg:  "操作失败",
		}, nil
	}
	return &proto.DouyinFavoriteActionResponse{
		StatusCode: 0,
		StatusMsg:  "操作成功",
	}, nil
}
