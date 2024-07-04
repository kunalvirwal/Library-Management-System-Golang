package models

type PENDING_REQUESTS struct {
	UUID  uint  `gorm:"not null"`
	BUID  uint  `gorm:"not null"`
	USER  USER  `gorm:"foreignKey:UUID;references:UUID"`
	BOOKS BOOKS `gorm:"foreignKey:BUID;references:BUID;constraint:OnDelete:CASCADE"`
	TYPE  bool  `gorm:"not null"`
}
