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
	IsFavorite    bool      `gorm:"default:false"`
}
