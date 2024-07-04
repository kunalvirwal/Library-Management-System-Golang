package models

import "time"

type BORROWING_HISTORY struct {
	UUID          uint       `gorm:"not null"`
	BUID          uint       `gorm:"not null"`
	USER          USER       `gorm:"foreignKey:UUID;references:UUID"`
	BOOKS         BOOKS      `gorm:"foreignKey:BUID;references:BUID"`
	CHECKOUT_DATE time.Time  `gorm:"not null"`
	CHECKIN_DATE  *time.Time `gorm:"default:null"`
}
