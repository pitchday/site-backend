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
	ReferralCode string
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

func (c *Contributor) GetByServiceId(serviceId string) (err error, isNotFound bool) {
	err = db.Raw("SELECT * FROM contributors WHERE id = (SELECT contributor_id FROM user_ims WHERE service_id = ?)", serviceId).Scan(&c).Error
	if err != nil && isNotFoundDBError(err) {
		isNotFound = true
	}

	return
}

func (c *Contributor) GetByReferralCode(referralCode string) (err error, isNotFound bool) {
	err = db.Raw("SELECT * FROM contributors WHERE referral_code = ?", referralCode).Scan(&c).Error
	if err != nil && isNotFoundDBError(err) {
		isNotFound = true
	}

	return
}

func (c *Contributor) Create() (err error, isDuplicated bool) {
	c.ReferralCode = uuid.NewV4().String()

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

func (c *ContributorList) Get() (err error) {
	err = db.Table("contributors").Find(&c).Error
	return
}

func MakeContributorFromTelegram(u tgbotapi.User, isMember bool, referredByCode string, privateChatId int64) bool {
	imageUrl, err := getUserAvatar(u)
	if err != nil {
		return false
	}

	referrer := Contributor{}

	if len(referredByCode) > 1 {
		err, isNotFound := referrer.GetByReferralCode(referredByCode)
		if err != nil && !isNotFound {
			Logger.Println("There was an error registering user with referral code", referredByCode)
		}
	}

	tx := db.Begin()

	contributor := Contributor{
		Name:         fmt.Sprintf("%s %s", u.FirstName, u.LastName),
		AvatarUrl:    imageUrl,
		ReferralCode: uuid.NewV4().String(),
		Description:  "member",
		ReferredBy:   referrer.Id,
	}

	if !isMember {
		contributor.DeletedAt.Scan(time.Now())
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

	if privateChatId != 0 {
		contributorIm.PrivateChat = privateChatId
	}

	err = tx.Create(&contributorIm).Error
	if err != nil {
		if isDuplicatedDBError(err) && isMember {
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
