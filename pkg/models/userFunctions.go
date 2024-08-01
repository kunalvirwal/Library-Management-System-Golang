package models

import (
	"github.com/kunalvirwal/go-mvc/pkg/utils"
	"gorm.io/gorm"
)

func SearchUserEmail(DB *gorm.DB, email string) (USER, bool) {
	var user USER
	res := DB.Where("EMAIL=?", email).First(&user)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {

			return user, false

		} else {
			utils.CheckNilErr(res.Error, "Error with DB while finding admin")
		}
	}
	return user, true
}

func SearchUserUUID(DB *gorm.DB, uuid int) (USER, bool) {
	var user USER
	res := DB.Where("UUID=?", uuid).First(&user)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return user, false
		} else {
			utils.CheckNilErr(res.Error, "Error with DB while finding admin")
		}
	}
	return user, true
}

func AdminExist(DB *gorm.DB) bool {
	var user USER
	res := DB.Where("ROLE=?", "admin").First(&user)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return false
		} else {
			utils.CheckNilErr(res.Error, "Error with DB while finding admin")
		}
	}
	return true
}

func CreateNewUser(DB *gorm.DB, name string, email string, phn_no int64, password string) {
	falsevalue := false
	var status *bool = &falsevalue
	user := USER{
		NAME:          name,
		EMAIL:         email,
		PHN_NO:        phn_no,
		PASSWORD:      password,
		ROLE:          "user",
		ADMIN_REQUEST: status,
	}
	result := DB.Create(&user)
	utils.CheckNilErr(result.Error, "Unable to create user")
}

func GetAllUsers(DB *gorm.DB) []USER { //TODO: can be used in users functionality also
	var dataset []USER
	DB.Find(&dataset)
	return dataset
}

func UpdateUserData(DB *gorm.DB, uuid int, name string, phn_no int64) {
	var user USER
	res := DB.Where("UUID=?", uuid).First(&user)
	utils.CheckNilErr(res.Error, "Error with DB while finding user")
	user.NAME = name
	user.PHN_NO = phn_no
	DB.Save(&user)
}

func GetAdminRequests(DB *gorm.DB) []USER {
	var dataset []USER
	DB.Where("ADMIN_REQUEST", true).Find(&dataset)
	return dataset
}

func CreateAdminReq(DB *gorm.DB, uuid int) {
	var user USER
	truevalue := true
	res := DB.Where("UUID=?", uuid).First(&user)
	utils.CheckNilErr(res.Error, "Error with DB while finding user")
	user.ADMIN_REQUEST = &truevalue
	DB.Save(&user)
}

func SetAdminReq(DB *gorm.DB, uuid int, status bool) {
	var user USER
	falsevalue := false
	res := DB.Where("UUID=?", uuid).First(&user)
	utils.CheckNilErr(res.Error, "Error with DB while finding user")
	if status {
		user.ADMIN_REQUEST = nil
		user.ROLE = "admin"
	} else {
		user.ADMIN_REQUEST = &falsevalue
	}
	DB.Save(&user)
}
