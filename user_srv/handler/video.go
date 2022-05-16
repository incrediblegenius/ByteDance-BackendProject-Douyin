package handler

import (
	"Douyin/proto"
	"context"
)

type VideosServer struct {
	proto.UnimplementedVideosServer
}

func (s *VideosServer) FavoriteAction(ctx context.Context, req *proto.DouyinFavoriteActionRequest) (*proto.DouyinFavoriteActionResponse, error) {
	return nil, nil
}
