package models

import "time"

type BORROWING_HISTORY struct {
	UUID          int        `gorm:"not null"`
	BUID          int        `gorm:"not null"`
	USER          USER       `gorm:"foreignKey:UUID;references:UUID"`
	BOOKS         BOOKS      `gorm:"foreignKey:BUID;references:BUID"`
	CHECKOUT_DATE time.Time  `gorm:"not null"`
	CHECKIN_DATE  *time.Time `gorm:"default:null"`
}
