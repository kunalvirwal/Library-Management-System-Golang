package models

type PENDING_REQUESTS struct {
	UUID  int   `gorm:"not null"`
	BUID  int   `gorm:"not null"`
	USER  USER  `gorm:"foreignKey:UUID;references:UUID"`
	BOOKS BOOKS `gorm:"foreignKey:BUID;references:BUID;constraint:OnDelete:CASCADE"`
	TYPE  bool  `gorm:"not null"`
}
