package models

import (
	"github.com/cheviz/pitchdayBackend/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"os"
	"fmt"
)

var db *gorm.DB
var Logger = log.New(os.Stdout, " ", log.Ldate|log.Ltime|log.Lshortfile)

// Setup initializes the Conn object
func Setup() error {
	var err error
	db, err = gorm.Open("mysql", config.Conf.DBString)
	if err != nil {
		return err
	}
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.LogMode(true)

	db.AutoMigrate(&NewsletterSubscription{})
	return err
}

func isDuplicatedDBError(err error) (isDuplicated bool) {
	if fmt.Sprint(err)[:27] == "Error 1062: Duplicate entry" {
		isDuplicated = true
	}
	return
}