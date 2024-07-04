package models

type USER struct {
	UUID          int    `gorm:"primaryKey;autoIncrement"`
	NAME          string `gorm:"type:varchar(50);not null"`
	EMAIL         string `gorm:"type:varchar(50);unique;not null"`
	PHN_NO        int64  `gorm:"type:bigint"`
	PASSWORD      string `gorm:"type:varchar(1000)"`
	ROLE          string `gorm:"type:varchar(10)"`
	ADMIN_REQUEST *bool  `gorm:"type:boolean;default:false"`
}
