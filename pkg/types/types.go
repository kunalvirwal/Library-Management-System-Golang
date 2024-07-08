package types

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UUID  int    `json:"uuid"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
	jwt.StandardClaims
}

type ContextKeyType string

type BooksCheckedOut struct {
	BUID          int
	NAME          string
	CHECKOUT_DATE time.Time
	Req           bool `gorm:"default:false"`
}

type BooksReturned struct {
	BUID          int
	NAME          string
	CHECKOUT_DATE time.Time
	CHECKIN_DATE  *time.Time `gorm:"not null"`
}

type PendingRequestData struct {
	UUID      int
	BUID      int
	USER_NAME string
	BOOK_NAME string
	TYPE      bool
}
