package handler

import (
	"UserServer/global"
	"UserServer/model"
	"UserServer/proto"
	"context"
	"sort"
	"time"
)

func (s *Server) GetUserFeed(ctx context.Context, req *proto.DouyinFeedRequest) (*proto.DouyinFeedResponse, error) {
	uid := 0
	if req.Token != "" {
		claim, err := global.Jwt.ParseToken(req.Token)
		if err != nil {
			return &proto.DouyinFeedResponse{
				StatusCode: -2,
				StatusMsg:  "token鉴权失败",
				VideoList:  []*proto.Video{&proto.Video{}}}, nil
		} else {
			uid = int(claim.Id)
		}
	}

	var videos []model.Video
	result := global.DB.Limit(30).Order("update_time desc").Find(&videos, "update_time < ?", time.UnixMilli(req.LatestTime))
	if result.Error != nil {
		// fmt.Println("查询失败")
		return &proto.DouyinFeedResponse{
			StatusCode: -1,
			StatusMsg:  "查询失败",
			VideoList:  []*proto.Video{&proto.Video{}},
		}, nil
	}
	if len(videos) == 0 {
		return &proto.DouyinFeedResponse{
			StatusCode: -2,
			StatusMsg:  "无更多视频",
			VideoList:  []*proto.Video{&proto.Video{}},
		}, nil
	}
	var vis []*proto.Video
	var nextTime int64
	if len(videos) > 0 {
		sort.Slice(videos, func(i, j int) bool {
			return videos[i].UpdatedAt.UnixMilli() > videos[j].UpdatedAt.UnixMilli()
		})
		nextTime = videos[len(videos)-1].UpdatedAt.UnixMilli()
	} else {
		nextTime = time.Now().UnixMilli()
	}
	for _, v := range videos {
		user, err := s.GetUserById(context.Background(), &proto.IdRequest{Id: int64(v.AuthorID), NeedToken: false})
		if err != nil {
			return nil, err
		}
		flag := false
		if uid != 0 {
			result := global.DB.First(&model.FavoriteVideo{}, "user_id = ? and video_id = ?", uid, v.ID)
			if result.RowsAffected != 0 {
				flag = true
			} else {
				flag = false
			}
		}
		vis = append(vis, &proto.Video{
			Id:            int64(v.ID),
			Author:        user,
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: int64(v.FavoriteCount),
			CommentCount:  int64(v.CommentCount),
			IsFavorite:    flag, // TODO 判断这个视频是否自己喜欢
		})
	}

	return &proto.DouyinFeedResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		VideoList:  vis,
		NextTime:   nextTime,
	}, nil
}
