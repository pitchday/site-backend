package models

import (
	"fmt"
	"github.com/badoux/checkmail"
	"github.com/go-telegram-bot-api/telegram-bot-api"
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
	DeletedAt   *time.Time `json:"-"`
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

func DeleteContributorUsingTelegramId(id int) (err error) {
	err = db.Exec("UPDATE contributors SET deleted_at = NOW() WHERE id IN (SELECT contributor_id FROM user_ims WHERE service_id = ?);", id).Error
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

func MakeContributorFromTelegram(u tgbotapi.User) bool {
	imageUrl, err := getUserAvatar(u)
	if err != nil {
		return false
	}

	tx := db.Begin()

	contributor := Contributor{
		Name:        fmt.Sprintf("%s %s", u.FirstName, u.LastName),
		AvatarUrl:   imageUrl,
		Description: "Contributor",
	}

	err = tx.Create(&contributor).Error
	if err != nil {
		tx.Rollback()
		return false
	}

	contributorIm := UserIm{
		ContributorId: contributor.Id,
		ServiceId:     fmt.Sprint(u.ID),
		Provider:      "Telegram",
	}
	err = tx.Create(&contributorIm).Error
	if err != nil {
		if isDuplicatedDBError(err) {
			tx.Rollback()

			err = db.Exec("UPDATE contributors SET deleted_at = NULL WHERE id IN (SELECT contributor_id FROM user_ims WHERE service_id = ?);", u.ID).Error
			if err != nil {
				return false
			}
			return true
		}

		tx.Rollback()
		return false
	}

	tx.Commit()
	return true
}

func getUserAvatar(user tgbotapi.User) (imageUrl string, err error) {
	conf := tgbotapi.UserProfilePhotosConfig{
		UserID: user.ID,
		Limit:  1,
	}

	avatars, err := telegramBot.GetUserProfilePhotos(conf)
	if err != nil {
		Logger.Println("Error retrieving user avatars,", err)
		return
	}

	avatarUrl, err := telegramBot.GetFileDirectURL(avatars.Photos[0][1].FileID)
	if err != nil {
		Logger.Println("Error getting user avatar url,", err)
		return
	}

	imageUrl, err = ImportImage(avatarUrl, fmt.Sprint(user.ID))

	return
}
