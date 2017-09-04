package models

import (
	"github.com/badoux/checkmail"
	"time"
)

type NewsletterSubscription struct {
	Id        int64  `json:"-"`
	Email     string `json:"email" gorm:"type:varchar(255);not null;unique"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (ns *NewsletterSubscription) Create() (err error, isDuplicated bool) {
	err = db.Create(&ns).Error
	if err != nil && isDuplicatedDBError(err) {
		isDuplicated = true
	}
	return
}

func (ns *NewsletterSubscription) Validate() (err error) {
	err = checkmail.ValidateFormat(ns.Email)
	return
}
