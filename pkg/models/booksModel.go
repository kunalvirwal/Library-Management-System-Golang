package models

type BOOKS struct {
	BUID        int    `gorm:"primaryKey;autoIncrement;unique"`
	NAME        string `gorm:"type:varchar(50);not null"`
	DESCRIPTION string `gorm:"type:varchar(2000)"`
	TOTAL       int    `gorm:"not null"`
	CHECKIN     int    `gorm:"check:TOTAL>=CHECKIN;not null"`
}
