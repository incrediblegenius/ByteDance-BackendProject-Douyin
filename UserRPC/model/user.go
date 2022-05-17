package model

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID        int            `gorm:"primarykey;type:int" json:"id"`
	CreatedAt time.Time      `gorm:"column:add_time" json:"-"`
	UpdatedAt time.Time      `gorm:"column:update_time" json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
	IsDeleted bool           `json:"-"`
}
type User struct {
	BaseModel
	UserName string `gorm:"index:idx_username,unique;type:varchar(40);not null"`
	Password string `gorm:"type:varchar(40);not null"`
}

type Relation struct {
	BaseModel
	UserFrom   User `gorm:"foreignkey:FollowFrom"`
	UserTo     User `gorm:"foreignkey:FollowTo"`
	FollowFrom int  `gorm:"index:idx_follow_from_to,unique;type:int;not null"`
	FollowTo   int  `gorm:"index:idx_follow_from_to,unique;index:idx_follow_to;type:int;not null"`
}
