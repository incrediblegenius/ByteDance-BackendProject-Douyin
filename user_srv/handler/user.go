package handler

import (
	"Douyin/model"
	"Douyin/proto"
	"Douyin/user_srv/global"
	"Douyin/user_srv/middleware"
	"context"
	"errors"
	"fmt"
	"sort"
	"time"
)

type UserRegisterServer struct {
	proto.UnimplementedUserRegisterServer
}

func (s *UserRegisterServer) Register(ctx context.Context, req *proto.DouyinUserRegisterRequest) (*proto.DouyinUserRegisterResponse, error) {
	username, password := req.Username, req.Password
	var user model.User
	result := global.DB.Where(&model.User{UserName: username}).First(&user)
	if result.RowsAffected == 1 {
		return &proto.DouyinUserRegisterResponse{
			StatusCode: 1,
			StatusMsg:  "用户名已存在",
		}, nil
	}
	user.UserName = username
	user.Password = password
	result = global.DB.Create(&user)
	if result.Error != nil {
		return &proto.DouyinUserRegisterResponse{
			StatusCode: 2,
			StatusMsg:  "注册失败",
		}, nil
	}
	token, _ := global.Jwt.CreateToken(middleware.CustomClaims{
		Id: int64(user.ID),
	})
	return &proto.DouyinUserRegisterResponse{
		StatusCode: 0,
		StatusMsg:  "注册成功",
		UserId:     int64(user.ID),
		Token:      token,
	}, nil
}

func (s *UserRegisterServer) Login(ctx context.Context, req *proto.DouyinUserRegisterRequest) (*proto.DouyinUserRegisterResponse, error) {
	username, password := req.Username, req.Password
	var user model.User
	result := global.DB.Where(&model.User{UserName: username}).First(&user)
	if result.RowsAffected == 0 {
		return &proto.DouyinUserRegisterResponse{
			StatusCode: 1,
			StatusMsg:  "用户不存在",
		}, nil
	}
	if user.Password != password {
		return &proto.DouyinUserRegisterResponse{
			StatusCode: 2,
			StatusMsg:  "用户密码错误",
		}, nil
	}
	token, _ := global.Jwt.CreateToken(middleware.CustomClaims{
		Id: int64(user.ID),
	})
	return &proto.DouyinUserRegisterResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		UserId:     int64(user.ID),
		Token:      token,
	}, nil
}

func (s *UserRegisterServer) GetUserById(ctx context.Context, req *proto.IdRequest) (*proto.User, error) {
	if req.NeedToken {
		claim, err := global.Jwt.ParseToken(req.Token)
		if err != nil {
			return nil, errors.New("token解析失败")
		} else if claim.Id != int64(req.Id) {
			return nil, errors.New("非法token")
		}
	}
	var user model.User
	ans := &proto.User{}
	result := global.DB.First(&user, req.Id)
	if result.RowsAffected == 0 {
		return nil, errors.New("未找到用户")
	}
	ans.Name = user.UserName
	ans.Id = int64(user.ID)
	var cnt int64
	global.DB.Where(&model.Relation{FollowFrom: int(req.Id)}).Count(&cnt)
	ans.FollowCount = cnt
	global.DB.Where(&model.Relation{FollowTo: int(req.Id)}).Count(&cnt)
	ans.FollowerCount = cnt
	ans.IsFollow = false
	return ans, nil
}

func (s *UserRegisterServer) GetUserFeed(ctx context.Context, req *proto.DouyinFeedRequest) (*proto.DouyinFeedResponse, error) {
	if req.Token != "" {
		_, err := global.Jwt.ParseToken(req.Token)
		if err != nil {
			return &proto.DouyinFeedResponse{
				StatusCode: -2,
				StatusMsg:  "token鉴权失败",
			}, nil
		}
	}
	var videos []model.Video
	result := global.DB.Limit(30).Find(&videos, "update_time > ?", time.UnixMilli(req.LatestTime))
	if result.Error != nil {
		// fmt.Println("查询失败")
		return &proto.DouyinFeedResponse{
			StatusCode: -1,
			StatusMsg:  "查询失败",
		}, nil
	}
	var vis []*proto.Video
	for _, v := range videos {
		user, err := s.GetUserById(context.Background(), &proto.IdRequest{Id: int64(v.AuthorID), NeedToken: false})
		if err != nil {
			fmt.Println(err.Error())
		}
		vis = append(vis, &proto.Video{
			Id:            int64(v.ID),
			Author:        user,
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: int64(v.FavoriteCount),
			CommentCount:  int64(v.CommentCount),
			IsFavorite:    v.IsFavorite,
		})
	}
	var nextTime int64
	if len(videos) > 0 {
		sort.Slice(videos, func(i, j int) bool {
			return videos[i].UpdatedAt.UnixMilli() < videos[j].UpdatedAt.UnixMilli()
		})
		nextTime = videos[len(videos)-1].UpdatedAt.UnixMilli()
	} else {
		nextTime = time.Now().UnixMilli()
	}
	return &proto.DouyinFeedResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		VideoList:  vis,
		NextTime:   nextTime,
	}, nil
}
