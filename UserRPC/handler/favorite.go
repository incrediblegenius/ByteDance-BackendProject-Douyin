package handler

import (
	"UserServer/global"
	"UserServer/model"
	"UserServer/proto"
	"context"

	"gorm.io/gorm"
)

type Server struct {
	proto.UnimplementedServerServer
}

func (s *Server) FavoriteAction(ctx context.Context, req *proto.DouyinFavoriteActionRequest) (*proto.DouyinFavoriteActionResponse, error) {
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
			StatusMsg:  "操作请求与实际不符",
		}, nil
	}
	return &proto.DouyinFavoriteActionResponse{
		StatusCode: 0,
		StatusMsg:  "操作成功",
	}, nil
}

func (s *Server) FavoriteList(ctx context.Context, req *proto.DouyinFavoriteListRequest) (*proto.DouyinFavoriteListResponse, error) {
	claim, err := global.Jwt.ParseToken(req.Token)
	if err != nil {
		return &proto.DouyinFavoriteListResponse{
			StatusCode: -2,
			StatusMsg:  "token鉴权失败",
			VideoList:  []*proto.Video{&proto.Video{}},
			// 返回nil前端报错
		}, nil
	}
	uid := req.UserId
	if uid == 0 {
		uid = claim.Id
	}
	var videoList []*model.FavoriteVideo
	global.DB.Where("user_id = ?", uid).Find(&videoList)
	if len(videoList) == 0 {
		return &proto.DouyinFavoriteListResponse{
			StatusCode: 0,
			StatusMsg:  "没有收藏视频",
			VideoList:  []*proto.Video{&proto.Video{}},
			// 返回nil前端报错
		}, nil
	}
	vis := make([]*proto.Video, len(videoList))
	for i := range videoList {
		vis[i], _ = s.GetVideoById(context.Background(), &proto.VideoIdRequest{
			VideoId:  int64(videoList[i].VideoID),
			SearchId: int64(uid),
		})
	}
	return &proto.DouyinFavoriteListResponse{
		StatusCode: 0,
		StatusMsg:  "操作成功",
		VideoList:  vis,
	}, nil
}

func (s *Server) GetVideoById(ctx context.Context, in *proto.VideoIdRequest) (*proto.Video, error) {
	video := model.Video{}
	result := global.DB.First(&video, "id = ?", in.VideoId)
	if result.Error != nil {
		return &proto.Video{}, result.Error
	}
	var author model.User
	result = global.DB.First(&author, "id = ?", video.AuthorID)
	if result.Error != nil {
		return &proto.Video{}, result.Error
	}
	var likeAuthor model.Relation
	var likeVideo model.FavoriteVideo
	r1 := global.DB.First(&likeAuthor, "follow_from = ? and follow_to = ?", in.SearchId, video.AuthorID)
	r2 := global.DB.First(&likeVideo, "user_id = ? and video_id = ?", in.SearchId, video.ID)
	return &proto.Video{
		Id: int64(video.ID),
		Author: &proto.User{
			Id:            int64(author.ID),
			Name:          author.UserName,
			FollowCount:   int64(author.FollowingCount),
			FollowerCount: int64(author.FollowerCount),
			IsFollow:      r1.RowsAffected != 0,
		},
		PlayUrl:       video.PlayUrl,
		CoverUrl:      video.CoverUrl,
		FavoriteCount: int64(video.FavoriteCount),
		CommentCount:  int64(video.CommentCount),
		IsFavorite:    r2.RowsAffected != 0,
	}, nil

}
