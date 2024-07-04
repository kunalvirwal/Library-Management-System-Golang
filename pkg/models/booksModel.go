package models

type BOOKS struct {
	BUID        uint   `gorm:"primaryKey;autoIncrement;unique"`
	NAME        string `gorm:"type:varchar(50);not null"`
	DESCRIPTION string `gorm:"type:varchar(2000)"`
	TOTAL       uint   `gorm:"not null"`
	CHECKIN     uint   `gorm:"check:TOTAL>=CHECKIN;not null"`
}
