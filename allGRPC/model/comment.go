package model

type Comment struct {
	BaseModel
	Video    Video  `gorm:"foreignkey:VideoID"`
	VideoID  int    `gorm:"index:idx_videoid;type:int;not null"`
	User     User   `gorm:"foreignkey:UserID"`
	UserID   int    `gorm:"index:idx_userid;type:int;not null"`
	Content  string `gorm:"type:varchar(255);not null"`
}
