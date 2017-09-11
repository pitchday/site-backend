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

func (a *Announcements) Get() (err error) {
	err = db.Table("announcements").Find(&a).Error
	return
}
