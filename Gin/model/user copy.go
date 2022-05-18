package model

import (
	"time"

	"gorm.io/gorm"
)

type BaseModeltest struct {
	ID        int            `gorm:"primarykey;type:int" json:"id"`
	CreatedAt time.Time      `gorm:"column:add_time" json:"-"`
	UpdatedAt time.Time      `gorm:"column:update_time" json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
	IsDeleted bool           `json:"-"`
}
type Usertest struct {
	BaseModel
	UserName  string      `gorm:"index:idx_username,unique;type:varchar(40);not null"`
	Password  string      `gorm:"type:varchar(40);not null"`
	Following *[]Usertest `gorm:"many2many:usertest_following;"`
	Follower  *[]Usertest `gorm:"many2many:user_follower;"`
}
