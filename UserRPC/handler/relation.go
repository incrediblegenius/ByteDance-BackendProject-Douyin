package handler

import (
	"UserServer/global"
	"UserServer/model"
	"UserServer/proto"
	"context"

	"gorm.io/gorm"
)

func (s *Server) RelationAction(ctx context.Context, req *proto.DouyinRelationActionRequest) (*proto.DouyinRelationActionResponse, error) {
	claim, err := global.Jwt.ParseToken(req.Token)
	if err != nil {
		return &proto.DouyinRelationActionResponse{
			StatusCode: -2,
			StatusMsg:  "token鉴权失败",
		}, nil
	}
	uid := req.UserId
	if uid == 0 {
		uid = claim.Id
	}
	var relation model.Relation
	result := global.DB.Where("follow_form = ? AND follow_to=?", uid, req.ToUserId).First(&relation)
	if result.RowsAffected == 0 && req.ActionType == 1 {
		r := global.DB.Transaction(func(tx *gorm.DB) error {
			// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
			if err := tx.Create(&model.Relation{
				FollowFrom: int(uid),
				FollowTo:   int(req.ToUserId),
			}).Error; err != nil {
				// 返回任何错误都会回滚事务
				return err
			}
			if err := tx.Model(&model.User{}).Where("id = ?", uid).Update("following_count", gorm.Expr("following_count + ?", 1)).Error; err != nil {
				return err
			}
			if err := tx.Model(&model.User{}).Where("id = ?", req.ToUserId).Update("follower_count", gorm.Expr("follower_count + ?", 1)).Error; err != nil {
				return err
			}
			// 返回 nil 提交事务
			return nil
		})
		if r != nil {
			return &proto.DouyinRelationActionResponse{
				StatusCode: -1,
				StatusMsg:  r.Error(),
			}, nil
		} else {
			return &proto.DouyinRelationActionResponse{
				StatusCode: 0,
				StatusMsg:  "关注成功",
			}, nil
		}
	} else if result.RowsAffected != 0 && req.ActionType == 2 {
		r := global.DB.Transaction(func(tx *gorm.DB) error {
			// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
			if err := tx.Delete(&relation).Error; err != nil {
				// 返回任何错误都会回滚事务
				return err
			}
			if err := tx.Model(&model.User{}).Where("id = ?", uid).Update("following_count", gorm.Expr("following_count - ?", 1)).Error; err != nil {
				return err
			}
			if err := tx.Model(&model.User{}).Where("id = ?", req.ToUserId).Update("follower_count", gorm.Expr("follower_count - ?", 1)).Error; err != nil {
				return err
			}
			// 返回 nil 提交事务
			return nil
		})
		if r != nil {
			return &proto.DouyinRelationActionResponse{
				StatusCode: -1,
				StatusMsg:  r.Error(),
			}, nil
		} else {
			return &proto.DouyinRelationActionResponse{
				StatusCode: 0,
				StatusMsg:  "取关成功",
			}, nil
		}
	}
	return &proto.DouyinRelationActionResponse{
		StatusCode: -1,
		StatusMsg:  "传参数错误",
	}, nil
}

func (s *Server) RelationFollowList(ctx context.Context, req *proto.DouyinRelationFollowListRequest) (*proto.DouyinRelationFollowListResponse, error) {
	_, err := global.Jwt.ParseToken(req.Token)
	if err != nil {
		return &proto.DouyinRelationFollowListResponse{
			StatusCode: -2,
			StatusMsg:  "token鉴权失败",
		}, nil
	}
	var relations []model.Relation
	result := global.DB.Where("follow_from = ?", req.UserId).Find(&relations)
	if result.RowsAffected == 0 {
		return &proto.DouyinRelationFollowListResponse{
			StatusCode: 0,
			StatusMsg:  "没有关注的人",
			UserList:   []*proto.User{},
		}, nil
	}
	userList := make([]*proto.User, len(relations))
	for i, v := range relations {
		userList[i], _ = s.GetUserById(ctx, &proto.IdRequest{
			Id:        int64(v.FollowTo),
			NeedToken: false,
			SearchId:  req.UserId,
		})
	}
	return &proto.DouyinRelationFollowListResponse{
		StatusCode: 0,
		StatusMsg:  "获取关注列表成功",
		UserList:   userList,
	}, nil
}

func (s *Server) RelationFollowerList(ctx context.Context, req *proto.DouyinRelationFollowerListRequest) (*proto.DouyinRelationFollowerListResponse, error) {
	_, err := global.Jwt.ParseToken(req.Token)
	if err != nil {
		return &proto.DouyinRelationFollowerListResponse{
			StatusCode: -2,
			StatusMsg:  "token鉴权失败",
		}, nil
	}
	var relations []model.Relation
	result := global.DB.Where("follow_to = ?", req.UserId).Find(&relations)
	if result.RowsAffected == 0 {
		return &proto.DouyinRelationFollowerListResponse{
			StatusCode: 0,
			StatusMsg:  "没有跟随者",
			UserList:   []*proto.User{},
		}, nil
	}
	userList := make([]*proto.User, len(relations))
	for i, v := range relations {
		userList[i], _ = s.GetUserById(ctx, &proto.IdRequest{
			Id:        int64(v.FollowFrom),
			NeedToken: false,
			SearchId:  req.UserId,
		})
	}
	return &proto.DouyinRelationFollowerListResponse{
		StatusCode: 0,
		StatusMsg:  "获取关注列表成功",
		UserList:   userList,
	}, nil
}
