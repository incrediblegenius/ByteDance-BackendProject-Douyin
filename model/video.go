package model

import "time"

type Video struct {
	BaseModel
	UpdatedAt     time.Time `gorm:"column:update_time;not null;index:idx_update" `
	Author        User      `gorm:"foreignkey:AuthorID"`
	AuthorID      int       `gorm:"index:idx_authorid;not null"`
	PlayUrl       string    `gorm:"type:varchar(255);not null"`
	CoverUrl      string    `gorm:"type:varchar(255)"`
	FavoriteCount int       `gorm:"default:0"`
	CommentCount  int       `gorm:"default:0"`
}

type FavoriteVideo struct {
	BaseModel
	Video   Video `gorm:"foreignkey:VideoID"`
	VideoID int   `gorm:"index:idx_videoid;not null"`
	User    User  `gorm:"foreignkey:UserID"`
	UserID  int   `gorm:"index:idx_userid;not null"`
}
