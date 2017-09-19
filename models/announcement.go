package models

import (
	"github.com/jinzhu/gorm"
)

type Announcement struct {
	gorm.Model
	MessageId    int   `json:"-" gorm:"index"`
	ChannelId    int64 `json:"-"`
	ResourceType string
	Resource     string `gorm:"size:4096"`
	Body         string `gorm:"size:4096"`
}

func (a *Announcement) Create() (err error) {
	err = db.Create(&a).Error
	return
}

func (a *Announcement) Edit() (err error) {
	err = db.Exec("UPDATE announcements SET body = ?, updated_at = NOW() WHERE message_id = ? AND channel_id = ?;", a.Body, a.MessageId, a.ChannelId).Error
	return
}

type Announcements []Announcement

func (a *Announcements) Get(limit, offset int) (err error) {
	if limit == 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	err = db.Table("announcements").Order("created_at DESC").Limit(limit).Offset(offset).Find(&a).Error
	return
}
