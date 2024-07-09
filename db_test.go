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
		PASSWORD:      "$2a$04$hgmWQ3ww4JHqjG5m37Aof.q8OmeDxoOqyKF2.mWBOewFeL2ZiGicS",
		EMAIL:         "admin@sdslabs.com",
		ROLE:          "admin",
		PHN_NO:        9999999999,
		UUID:          1,
		ADMIN_REQUEST: nil,
	}
	if !found {
		t.Errorf("Admin not found")
	}
	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}
}

func TestSearchUserUUID(t *testing.T) {
	db, err := models.Connection()
	utils.CheckNilErr(err, "Unable to create Db instance for test")

	got, found := models.SearchUserUUID(db, 1)
	want := models.USER{
		NAME:          "admin",
		PASSWORD:      "$2a$04$hgmWQ3ww4JHqjG5m37Aof.q8OmeDxoOqyKF2.mWBOewFeL2ZiGicS",
		EMAIL:         "admin@sdslabs.com",
		ROLE:          "admin",
		PHN_NO:        9999999999,
		UUID:          1,
		ADMIN_REQUEST: nil,
	}
	if !found {
		t.Errorf("Admin not found")
	}
	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}
}
