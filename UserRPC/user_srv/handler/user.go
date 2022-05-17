package handler

import (
	"UserServer/model"
	"UserServer/proto"
	"UserServer/user_srv/global"
	"UserServer/user_srv/middleware"
	"context"
	"errors"
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
	var r []model.Relation
	global.DB.Where(&model.Relation{FollowFrom: int(req.Id)}).Find(&r)
	ans.FollowCount = int64(len(r))
	global.DB.Where(&model.Relation{FollowTo: int(req.Id)}).Find(&r)
	ans.FollowerCount = int64(len(r))
	ans.IsFollow = false
	return ans, nil
}
