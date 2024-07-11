package models

import (
	"testing"

	"github.com/kunalvirwal/go-mvc/pkg/models"
	"github.com/kunalvirwal/go-mvc/pkg/utils"
)

func TestSearchUserEmail(t *testing.T) {
	db, err := models.Connection()
	utils.CheckNilErr(err, "Unable to create Db instance for test")

	got, found := models.SearchUserEmail(db, "admin@sdslabs.com")
	want := models.USER{
		NAME:          "admin",
		PASSWORD:      "A",
		EMAIL:         "admin@sdslabs.com",
		ROLE:          "admin",
		PHN_NO:        9999999999,
		UUID:          1,
		ADMIN_REQUEST: nil,
	}
	if !found {
		t.Errorf("Admin not found")
	}
	if got.NAME != want.NAME || got.EMAIL != want.EMAIL || got.UUID != want.UUID || got.PHN_NO != want.PHN_NO || got.ROLE != want.ROLE || !(utils.MatchHashtoPassword(got.PASSWORD, want.PASSWORD)) {
		t.Errorf("got %v, wanted %v", got, want)
	}
}

func TestSearchUserUUID(t *testing.T) {
	db, err := models.Connection()
	utils.CheckNilErr(err, "Unable to create Db instance for test")

	got, found := models.SearchUserUUID(db, 1)
	want := models.USER{
		NAME:          "admin",
		PASSWORD:      "A",
		EMAIL:         "admin@sdslabs.com",
		ROLE:          "admin",
		PHN_NO:        9999999999,
		UUID:          1,
		ADMIN_REQUEST: nil,
	}
	if !found {
		t.Errorf("Admin not found")
	}

	// using each condn indiviually because we have to comapre a hashed password to a string
	if got.NAME != want.NAME || got.EMAIL != want.EMAIL || got.UUID != want.UUID || got.PHN_NO != want.PHN_NO || got.ROLE != want.ROLE || !(utils.MatchHashtoPassword(got.PASSWORD, want.PASSWORD)) {
		t.Errorf("got %v, wanted %v", got, want)
	}
}
