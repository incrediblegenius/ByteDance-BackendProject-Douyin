package model

import "time"

type Videotest struct {
	BaseModel
	UpdatedAt    time.Time  `gorm:"column:update_time;not null;index:idx_update" `
	Author       User       `gorm:"foreignkey:AuthorID"`
	AuthorID     int        `gorm:"index:idx_authorid;not null"`
	PlayUrl      string     `gorm:"type:varchar(255);not null"`
	CoverUrl     string     `gorm:"type:varchar(255)"`
	Likers       []Usertest `gorm:"many2many:videotest_usertest;"`
	CommentCount int        `gorm:"default:0"`
}
