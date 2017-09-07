package models

import "github.com/jinzhu/gorm"

type UserIm struct {
	gorm.Model
	ServiceId     string `gorm:"unique"`
	Provider      string
	ContributorId int64 `gorm:"index"`
}
