package handler

import (
	"Douyin/cfg"
	"Douyin/model"
	"Douyin/proto/userproto"
	"Douyin/user_srv/global"
	"context"
	"fmt"
)

func (s *UserRegisterServer) PublishAction(ctx context.Context, req *userproto.DouyinPublishActionRequest) (*userproto.DouyinPublishActionResponse, error) {
	claim, err := global.Jwt.ParseToken(req.Token)
	if err != nil {
		// os.Remove("../videos/" + req.VideoName)
		return &userproto.DouyinPublishActionResponse{
			StatusCode: -2,
			StatusMsg:  "token鉴权失败",
		}, nil
	}
	filename := req.VideoName

	id := int(claim.Id)

	video := &model.Video{
		AuthorID: id,
		PlayUrl:  fmt.Sprintf("http://%s:%d/videos/%s.mp4", cfg.ServerIP, cfg.ServerPort, filename),
		CoverUrl: fmt.Sprintf("http://%s:%d/covers/%s.png", cfg.ServerIP, cfg.ServerPort, filename),
	}
	result := global.DB.Create(&video)
	if result.Error != nil {
		return &userproto.DouyinPublishActionResponse{
			StatusCode: -1,
			StatusMsg:  "上传失败",
		}, nil
	}
	return &userproto.DouyinPublishActionResponse{
		StatusCode: 0,
		StatusMsg:  "视频上传成功",
	}, nil
}

func (s *UserRegisterServer) PublishList(ctx context.Context, req *userproto.DouyinPublishListRequest) (*userproto.DouyinPublishListResponse, error) {
	claim, err := global.Jwt.ParseToken(req.Token)
	if err != nil {
		return &userproto.DouyinPublishListResponse{
			StatusCode: -2,
			StatusMsg:  "token鉴权失败",
			VideoList:  []*userproto.Video{&userproto.Video{}},
		}, nil
	}

	id := int(claim.Id)
	var videos []model.Video
	result := global.DB.Where("author_id = ?", id).Find(&videos)
	if result.Error != nil || len(videos) == 0 {
		return &userproto.DouyinPublishListResponse{
			StatusCode: -1,
			StatusMsg:  "获取视频列表失败",
			VideoList:  []*userproto.Video{&userproto.Video{}},
		}, nil
	}
	vs := make([]*userproto.Video, len(videos))
	user, _ := s.GetUserById(context.Background(), &userproto.IdRequest{
		Id:        int64(id),
		NeedToken: false,
	})
	flag := false
	for i, v := range videos {
		result := global.DB.First(&model.FavoriteVideo{}, "video_id = ? and user_id = ?", v.ID, id)
		if result.RowsAffected != 0 {
			flag = true
		}
		vs[i] = &userproto.Video{
			Id:            int64(v.ID),
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			IsFavorite:    flag,
			Author:        user,
			FavoriteCount: int64(v.FavoriteCount),
			CommentCount:  int64(v.CommentCount),
		}
	}
	// fmt.Println(vs)
	return &userproto.DouyinPublishListResponse{
		StatusCode: 0,
		StatusMsg:  "获取视频列表成功",
		VideoList:  vs,
	}, nil
}
