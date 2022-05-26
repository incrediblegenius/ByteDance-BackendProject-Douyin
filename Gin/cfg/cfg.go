package cfg

type ServiceConfig struct {
	UserSrv     string `json:"user_srv"`
	FeedSrv     string `json:"feed_srv"`
	PublishSrv  string `json:"publish_srv"`
	FavoriteSrv string `json:"favorite_srv"`
	RelationSrv string `json:"relation_srv"`
	CommentSrv  string `json:"comment_srv"`
}
