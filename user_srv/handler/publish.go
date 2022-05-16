package handler

import (
	"Douyin/cfg"
	"Douyin/model"
	"Douyin/proto"
	"Douyin/user_srv/global"
	"context"
	"fmt"
	"os"
)

func (s *UserRegisterServer) PublishAction(ctx context.Context, req *proto.DouyinPublishActionRequest) (*proto.DouyinPublishActionResponse, error) {
	claim, err := global.Jwt.ParseToken(req.Token)
	if err != nil {
		os.Remove("../videos/" + req.VideoName)
		return &proto.DouyinPublishActionResponse{
			StatusCode: -2,
			StatusMsg:  "token鉴权失败",
		}, nil
	}
	filename := req.VideoName

	id := int(claim.Id)

	video := &model.Video{
		AuthorID:   id,
		PlayUrl:    fmt.Sprintf("http://%s:%d/videos/%s", cfg.ServerIP, cfg.ServerPort, filename),
		CoverUrl:   fmt.Sprintf("http://%s:%d/covers/%s", cfg.ServerIP, cfg.ServerPort, filename),
		IsFavorite: false,
	}
	result := global.DB.Create(&video)
	if result.Error != nil {
		return &proto.DouyinPublishActionResponse{
			StatusCode: -1,
			StatusMsg:  "上传失败",
		}, nil
	}
	return &proto.DouyinPublishActionResponse{
		StatusCode: 0,
		StatusMsg:  "视频上传成功",
	}, nil
}

func (s *UserRegisterServer) PublishList(ctx context.Context, req *proto.DouyinPublishListRequest) (*proto.DouyinPublishListResponse, error) {
	claim, err := global.Jwt.ParseToken(req.Token)
	if err != nil {
		return &proto.DouyinPublishListResponse{
			StatusCode: -2,
			StatusMsg:  "token鉴权失败",
			VideoList:  []*proto.Video{},
		}, nil
	}
	// TODO

	id := int(claim.Id)
	var videos []model.Video
	result := global.DB.Where("author_id = ?", id).Find(&videos)
	if result.Error != nil {
		return &proto.DouyinPublishListResponse{
			StatusCode: -1,
			StatusMsg:  "获取视频列表失败",
			VideoList:  []*proto.Video{},
		}, nil
	}
	vs := make([]*proto.Video, len(videos))
	user, _ := s.GetUserById(context.Background(), &proto.IdRequest{
		Id:        int64(id),
		NeedToken: false,
	})
	for i, v := range videos {

		vs[i] = &proto.Video{
			Id:         int64(v.ID),
			PlayUrl:    v.PlayUrl,
			CoverUrl:   v.CoverUrl,
			IsFavorite: v.IsFavorite,
			Author:     user,
		}
	}
	return &proto.DouyinPublishListResponse{
		StatusCode: 0,
		StatusMsg:  "获取视频列表成功",
		VideoList:  vs,
	}, nil
}
