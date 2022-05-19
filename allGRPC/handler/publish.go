package handler

import (
	"UserServer/global"
	"UserServer/model"
	"UserServer/proto"
	"context"
	"fmt"
)

func (s *Server) PublishAction(ctx context.Context, req *proto.DouyinPublishActionRequest) (*proto.DouyinPublishActionResponse, error) {
	claim, err := global.Jwt.ParseToken(req.Token)
	if err != nil {
		// os.Remove("../videos/" + req.VideoName)
		return &proto.DouyinPublishActionResponse{
			StatusCode: -2,
			StatusMsg:  "token鉴权失败",
		}, nil
	}
	filename := req.VideoName

	id := int(claim.Id)

	video := &model.Video{
		AuthorID: id,
		PlayUrl:  fmt.Sprintf("http://%s:%d/videos/%s.mp4", global.ServerConfig.StaticInfo.Host, global.ServerConfig.StaticInfo.Port, filename),
		CoverUrl: fmt.Sprintf("http://%s:%d/covers/%s.png", global.ServerConfig.StaticInfo.Host, global.ServerConfig.StaticInfo.Port, filename),
		Title:    req.Title,
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

func (s *Server) PublishList(ctx context.Context, req *proto.DouyinPublishListRequest) (*proto.DouyinPublishListResponse, error) {
	claim, err := global.Jwt.ParseToken(req.Token)
	if err != nil {
		return &proto.DouyinPublishListResponse{
			StatusCode: -2,
			StatusMsg:  "token鉴权失败",
			VideoList:  []*proto.Video{&proto.Video{}},
			// 返回nil前端报错
		}, nil
	}
	var id int
	if req.UserId == 0 {
		id = int(claim.Id)
	} else {
		id = int(req.UserId)
	}
	var videos []model.Video
	result := global.DB.Where("author_id = ?", id).Find(&videos)
	if result.Error != nil || len(videos) == 0 {
		return &proto.DouyinPublishListResponse{
			StatusCode: -1,
			StatusMsg:  "获取视频列表失败",
			VideoList:  []*proto.Video{&proto.Video{}},
			// 返回nil前端报错
		}, nil
	}
	vs := make([]*proto.Video, len(videos))
	user, _ := s.GetUserById(context.Background(), &proto.IdRequest{
		Id:        int64(id),
		NeedToken: false,
	})
	flag := false
	for i, v := range videos {
		result := global.DB.First(&model.FavoriteVideo{}, "video_id = ? and user_id = ?", v.ID, id)
		if result.RowsAffected != 0 {
			flag = true
		}
		vs[i] = &proto.Video{
			Id:            int64(v.ID),
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			IsFavorite:    flag,
			Author:        user,
			FavoriteCount: int64(v.FavoriteCount),
			CommentCount:  int64(v.CommentCount),
			Title:         v.Title,
		}
	}
	// fmt.Println(vs)
	return &proto.DouyinPublishListResponse{
		StatusCode: 0,
		StatusMsg:  "获取视频列表成功",
		VideoList:  vs,
	}, nil
}
