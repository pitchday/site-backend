package models

import "github.com/jinzhu/gorm"

type UserIm struct {
	gorm.Model
	ServiceId     string `gorm:"unique"`
	Provider      string
	PrivateChat   int64
	ContributorId int64 `gorm:"index"`
}
