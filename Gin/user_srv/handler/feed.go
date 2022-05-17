package handler

import (
	"Douyin/model"
	"Douyin/proto/userproto"
	"Douyin/user_srv/global"
	"context"
	"sort"
	"time"
)

func (s *UserRegisterServer) GetUserFeed(ctx context.Context, req *userproto.DouyinFeedRequest) (*userproto.DouyinFeedResponse, error) {
	uid := 0
	if req.Token != "" {
		claim, err := global.Jwt.ParseToken(req.Token)
		if err != nil {
			return &userproto.DouyinFeedResponse{
				StatusCode: -2,
				StatusMsg:  "token鉴权失败",
				VideoList:  []*userproto.Video{&userproto.Video{}}}, nil
		} else {
			uid = int(claim.Id)
		}
	}

	var videos []model.Video
	result := global.DB.Limit(30).Order("update_time desc").Find(&videos, "update_time < ?", time.UnixMilli(req.LatestTime))
	if result.Error != nil {
		// fmt.Println("查询失败")
		return &userproto.DouyinFeedResponse{
			StatusCode: -1,
			StatusMsg:  "查询失败",
			VideoList:  []*userproto.Video{&userproto.Video{}},
		}, nil
	}
	if len(videos) == 0 {
		return &userproto.DouyinFeedResponse{
			StatusCode: -2,
			StatusMsg:  "无更多视频",
			VideoList:  []*userproto.Video{&userproto.Video{}},
		}, nil
	}
	var vis []*userproto.Video
	for _, v := range videos {
		user, err := s.GetUserById(context.Background(), &userproto.IdRequest{Id: int64(v.AuthorID), NeedToken: false})
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
		vis = append(vis, &userproto.Video{
			Id:            int64(v.ID),
			Author:        user,
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: int64(v.FavoriteCount),
			CommentCount:  int64(v.CommentCount),
			IsFavorite:    flag, // TODO 判断这个视频是否自己喜欢
		})
	}
	var nextTime int64
	if len(videos) > 0 {
		sort.Slice(videos, func(i, j int) bool {
			return videos[i].UpdatedAt.UnixMilli() > videos[j].UpdatedAt.UnixMilli()
		})
		nextTime = videos[0].UpdatedAt.UnixMilli()
	} else {
		nextTime = time.Now().UnixMilli()
	}
	return &userproto.DouyinFeedResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		VideoList:  vis,
		NextTime:   nextTime,
	}, nil
}
