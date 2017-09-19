package models

import (
	"time"
	"github.com/jinzhu/gorm"
)

type ValueCache struct {
	gorm.Model
	Key       string
	Value     string
	ExpiresAt time.Time
}

func (vc *ValueCache) GetByKey() (err error) {
	err = db.Raw("SELECT * FROM value_caches vc WHERE vc.key = ?;", vc.Key).Scan(&vc).Error
	return
}

func (vc *ValueCache) Update() (err error) {
	err = db.Save(&vc).Error
	return
}