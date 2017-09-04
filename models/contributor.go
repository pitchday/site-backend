package models

import (
	"github.com/badoux/checkmail"
	"time"
)

type ContributorList []Contributor

type Contributor struct {
	Id          int64 `json:"-"`
	Name        string
	Description string
	Email       string
	Link        string
	AvatarUrl   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

func (c *Contributor) Get() (err error, isNotFound bool) {
	err = db.First(&c).Error
	if err != nil && isNotFoundDBError(err) {
		isNotFound = true
	}

	//Hide personal information
	c.Email = ""

	return
}

func (c *Contributor) Create() (err error, isDuplicated bool) {
	err = db.Create(&c).Error
	if err != nil && isDuplicatedDBError(err) {
		isDuplicated = true
	}
	return
}

func (c *Contributor) Delete() (err error) {
	err = db.Delete(&c).Error
	return
}

func (c *Contributor) Validate() (err error) {
	if len(c.Email) > 1 {
		err = checkmail.ValidateFormat(c.Email)
	}
	return
}


func (c *ContributorList) Get() (err error) {
	err = db.Table("contributors").Find(&c).Error
	return
}
