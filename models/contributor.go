package models

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/satori/go.uuid"
	"time"
)

type ContributorList []Contributor

type Contributor struct {
	Id           int64
	Name         string
	Description  string
	Link         string
	AvatarUrl    string
	ServiceId    int `gorm:"unique"`
	PrivateChat  int64
	ReferralCode string `gorm:"unique"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    mysql.NullTime `json:"-"`
	ReferredBy   int64          `json:"-"`
}

func (c *Contributor) Get() (err error, isNotFound bool) {
	err = db.First(&c).Error
	if err != nil && isNotFoundDBError(err) {
		isNotFound = true
	}

	return
}

func (c *Contributor) GetByServiceId(serviceId string) (err error) {
	err = db.Raw("SELECT * FROM contributors WHERE service_id = ?;", serviceId).Scan(&c).Error
	return
}

func (c *Contributor) GetByReferralCode(referralCode string) (err error) {
	err = db.Raw("SELECT * FROM contributors WHERE referral_code = ?;", referralCode).Scan(&c).Error
	return
}

func (c *Contributor) Create() (err error) {
	c.ReferralCode = uuid.NewV4().String()
	err = db.Create(&c).Error
	return
}

func (c *Contributor) Delete() (err error) {
	err = db.Delete(&c).Error
	return
}

func DeleteContributorUsingTelegramId(id int) (err error) {
	err = db.Exec("UPDATE contributors SET deleted_at = NOW() WHERE service_id = ?;", id).Error
	return
}

func (c *ContributorList) Get() (err error) {
	err = db.Table("contributors").Find(&c).Error
	return
}

func (c *Contributor) MakeMember() (err error) {
	c.DeletedAt.Valid = false
	err = db.Exec("UPDATE contributors SET deleted_at = null WHERE service_id = ?;", c.ServiceId).Error
	return
}

func MakeContributorFromTelegram(u tgbotapi.User, isMember bool, referredByCode string, privateChat int64) bool {
	imageUrl, err := getUserAvatar(u)
	if err != nil {
		return false
	}

	referrer := Contributor{}
	if len(referredByCode) > 1 {
		err := referrer.GetByReferralCode(referredByCode)
		if err != nil && !isNotFoundDBError(err) {
			Logger.Println("There was an error registering user with referral code", referredByCode)
		}
	}

	contributor := Contributor{
		Name:        fmt.Sprintf("%s %s", u.FirstName, u.LastName),
		AvatarUrl:   imageUrl,
		Description: "member",
		ServiceId:   u.ID,
		ReferredBy:  referrer.Id,
		PrivateChat: privateChat,
	}

	if !isMember {
		contributor.DeletedAt.Scan(time.Now())
	}

	err = contributor.Create()
	if err != nil {
		return false
	}
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

	if len(avatars.Photos) < 1 {
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
