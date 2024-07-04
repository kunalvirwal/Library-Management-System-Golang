package models

import (
	"github.com/kunalvirwal/go-mvc/pkg/utils"
	"gorm.io/gorm"
)

func AdminExist(DB *gorm.DB) bool {

	var user USER
	res := DB.Where("ROLE=?", "admin").First(&user)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return false
		} else {
			utils.CheckNilErr(res.Error, "Error with DB while finding admin") //// if an unknown error occurs here then panic will stop the control flow, if checkNilErr has to be removed then adequate return has too be added
		}
	}
	return true

}

func SearchUserEmail(DB *gorm.DB, email string) (USER, bool) {
	var user USER
	res := DB.Where("EMAIL=?", email).First(&user)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {

			return user, false

		} else {
			utils.CheckNilErr(res.Error, "Error with DB while finding admin") //// if an unknown error occurs here then panic will stop the control flow, if checkNilErr has to be removed then adequate return has too be added
		}
	}
	return user, true
}

func CreateNewUser(DB *gorm.DB, name string, email string, phn_no int64, password string) {
	user := USER{
		NAME:     name,
		EMAIL:    email,
		PHN_NO:   phn_no,
		PASSWORD: password,
		ROLE:     "user",
	}
	result := DB.Create(&user)
	utils.CheckNilErr(result.Error, "Unable to create user")
	// fmt.Println("Account created")
}
