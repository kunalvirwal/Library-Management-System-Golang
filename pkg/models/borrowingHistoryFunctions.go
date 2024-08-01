package models

import (
	"time"

	"github.com/kunalvirwal/go-mvc/pkg/utils"
	"gorm.io/gorm"
)

func IsCheckedOutByUser(DB *gorm.DB, buid int, uuid int) bool {
	var dataset BORROWING_HISTORY
	res := DB.Where(&BORROWING_HISTORY{UUID: uuid, BUID: buid}).Where("checkin_date IS NULL").First(&dataset)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return false
		} else {
			utils.CheckNilErr(res.Error, "Error with DB while finding book")
		}
	}
	return true
}

func CreateCheckoutEntry(DB *gorm.DB, buid int, uuid int) {
	record := BORROWING_HISTORY{
		UUID:          uuid,
		BUID:          buid,
		CHECKOUT_DATE: time.Now(),
		CHECKIN_DATE:  nil,
	}
	result := DB.Create(&record)
	utils.CheckNilErr(result.Error, "Unable to create recorf for book checkout")
}

func UpdateCheckinEntry(DB *gorm.DB, buid int, uuid int) {
	var record BORROWING_HISTORY
	time_now := time.Now()
	DB.Model(&record).Where(&BORROWING_HISTORY{UUID: uuid, BUID: buid}).Where("checkin_date IS NULL").Update("checkin_date", &time_now)
}
