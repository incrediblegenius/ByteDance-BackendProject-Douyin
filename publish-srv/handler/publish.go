package handler

import (
	"UserServer/global"
	"UserServer/model"
	"UserServer/proto"
	"context"
	"errors"
	"fmt"
)

type Server struct {
	proto.UnimplementedServerServer
}

func (s *Server) PublishAction(ctx context.Context, req *proto.DouyinPublishActionRequest) (*proto.DouyinPublishActionResponse, error) {
	claim, err := global.Jwt.ParseToken(req.Token)
	if err != nil {
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

	cid := claim.Id
	uid := req.UserId
	if uid == 0 {
		uid = cid // 如果没有传入userid，默认为自己
	}

	var videos []model.Video
	result := global.DB.Where("author_id = ?", uid).Find(&videos)
	if result.Error != nil {
		return &proto.DouyinPublishListResponse{
			StatusCode: -1,
			StatusMsg:  "获取视频列表失败",
			VideoList:  []*proto.Video{&proto.Video{}},
			// 返回nil前端报错
		}, nil
	}
	vs := make([]*proto.Video, len(videos))
	user, _ := GetUserById(&Request{
		Id: uid,
	})
	flag := false
	for i, v := range videos {

		result := global.DB.First(&model.FavoriteVideo{}, "video_id = ? and user_id = ?", v.ID, cid) // 是否喜欢是针对登陆用户（token）的
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
	return &proto.DouyinPublishListResponse{
		StatusCode: 0,
		StatusMsg:  "获取视频列表成功",
		VideoList:  vs,
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
