package model

import "time"

type Video struct {
	BaseModel
	CreatedAt     time.Time `gorm:"column:add_time;not null;index:idx_add" `
	Author        User      `gorm:"foreignkey:AuthorID"`
	AuthorID      int       `gorm:"index:idx_authorid;not null"`
	PlayUrl       string    `gorm:"type:varchar(255);not null"`
	CoverUrl      string    `gorm:"type:varchar(255)"`
	FavoriteCount int       `gorm:"default:0"`
	CommentCount  int       `gorm:"default:0"`
	Title         string    `gorm:"type:varchar(50);not null"`
}

type FavoriteVideo struct {
	BaseModel
	User    User  `gorm:"foreignkey:UserID"`
	UserID  int   `gorm:"index:idx_userid_videoid,unique;not null"`
	Video   Video `gorm:"foreignkey:VideoID"`
	VideoID int   `gorm:"index:idx_userid_videoid,unique;index:idx_videoid;not null"`
}
