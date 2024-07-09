package models

import (
	"time"

	"github.com/kunalvirwal/go-mvc/pkg/types"
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

func SearchUserUUID(DB *gorm.DB, uuid int) (USER, bool) {
	var user USER
	res := DB.Where("UUID=?", uuid).First(&user)
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
	falsevalue := false
	var status *bool = &falsevalue
	// if Admin_req_null { //used for create admin
	// 	status = nil
	// }
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
	// fmt.Println("Account created")
}

func GetCheckedOutBooks(DB *gorm.DB, uuid int) []types.BooksCheckedOut {
	var books []types.BooksCheckedOut
	DB.Table("borrowing_histories").Select("borrowing_histories.B_UID as BUID, borrowing_histories.CHECKOUT_DATE, books.NAME").Joins("join books on books.B_UID = borrowing_histories.B_UID ").Where(&BORROWING_HISTORY{UUID: uuid}).Where("checkin_date IS NULL").Scan(&books)
	return books

}

func GetAllPastCheckedInBooks(DB *gorm.DB, uuid int) []types.BooksReturned {
	var books []types.BooksReturned
	DB.Table("borrowing_histories").Select("borrowing_histories.B_UID as BUID, borrowing_histories.CHECKOUT_DATE,borrowing_histories.CHECKIN_DATE, books.NAME").Joins("join books on books.B_UID = borrowing_histories.B_UID ").Where(&BORROWING_HISTORY{UUID: uuid}).Where("checkin_date IS NOT NULL").Scan(&books)
	return books
}

func GetAllPendingCheckinReqByUUID(DB *gorm.DB, uuid int) []PENDING_REQUESTS {
	var reqs []PENDING_REQUESTS
	DB.Where(&PENDING_REQUESTS{UUID: uuid, TYPE: true}).Find(&reqs)
	return reqs
}

func GetAllBooks(DB *gorm.DB) []BOOKS {
	var dataset []BOOKS
	DB.Find(&dataset)
	return dataset
}

func GetAllUsers(DB *gorm.DB) []USER { ///////can be used in users functionality also
	var dataset []USER
	DB.Find(&dataset)
	return dataset
}

func GetAllPendingRequests(DB *gorm.DB) []PENDING_REQUESTS {
	var dataset []PENDING_REQUESTS
	DB.Find(&dataset)
	return dataset
}

func PendingReqExist(DB *gorm.DB, buid int, uuid int) (PENDING_REQUESTS, bool) {
	var req PENDING_REQUESTS
	res := DB.Where(&PENDING_REQUESTS{UUID: uuid, BUID: buid}).First(&req)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return req, false
		} else {
			utils.CheckNilErr(res.Error, "Error with DB while finding pending request") //// if an unknown error occurs here then panic will stop the control flow, if checkNilErr has to be removed then adequate return should be added
		}
	}

	return req, true

}

func GetBookCount(DB *gorm.DB) int {
	var dataset []BOOKS
	result := DB.First(&dataset)
	return int(result.RowsAffected)
}

func GetBook(DB *gorm.DB, buid int) (BOOKS, bool) {
	var book BOOKS
	res := DB.Where("B_UID=?", buid).First(&book)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return book, false
		} else {
			utils.CheckNilErr(res.Error, "Error with DB while finding book") //// if an unknown error occurs here then panic will stop the control flow, if checkNilErr has to be removed then adequate return should be added
		}
	}
	return book, true
}

func IsCheckedOutByUser(DB *gorm.DB, buid int, uuid int) bool {
	var dataset BORROWING_HISTORY
	res := DB.Where(&BORROWING_HISTORY{UUID: uuid, BUID: buid}).Where("checkin_date IS NULL").First(&dataset)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return false
		} else {
			utils.CheckNilErr(res.Error, "Error with DB while finding book") //// if an unknown error occurs here then panic will stop the control flow, if checkNilErr has to be removed then adequate return should be added
		}
	}
	return true
}

func CreateNewBookReq(DB *gorm.DB, buid int, uuid int, typ bool) {
	req := PENDING_REQUESTS{UUID: uuid, BUID: buid, TYPE: typ}
	res := DB.Create(req)
	if res.Error != nil {
		utils.CheckNilErr(res.Error, "Unable to create book Pending request ")
	}
}

func CreateCheckoutEntry(DB *gorm.DB, buid int, uuid int) {
	record := BORROWING_HISTORY{
		UUID:          uuid,
		BUID:          buid,
		CHECKOUT_DATE: time.Now(),
		CHECKIN_DATE:  nil,
	}
	result := DB.Create(&record)
	utils.CheckNilErr(result.Error, "Unable to create recorf for book checkout")

}

func UpdateCheckinEntry(DB *gorm.DB, buid int, uuid int) {
	var record BORROWING_HISTORY
	time_now := time.Now()
	DB.Model(&record).Where(&BORROWING_HISTORY{UUID: uuid, BUID: buid}).Where("checkin_date IS NULL").Update("checkin_date", &time_now)
}
func EditBook(DB *gorm.DB, buid int, name string, desc string, Checkin int, Total int) {
	book, found := GetBook(DB, buid)
	if found && name != "" && desc != "" && Checkin <= Total && Total >= 0 {
		book.NAME = name
		book.DESCRIPTION = desc
		book.CHECKIN = Checkin
		book.TOTAL = Total
		DB.Save(&book)

	} else {
		panic("Book does not exist to edit or invalid details")
	}
}

func GetDataofAllPendingRequests(DB *gorm.DB) []types.PendingRequestData {
	var records []types.PendingRequestData
	DB.Table("pending_requests").Select("users.UUID, books.B_UID as BUID, users.NAME as USER_NAME, books.NAME as BOOK_NAME, pending_requests.TYPE").Joins("join books on books.B_UID = pending_requests.B_UID").Joins("join users on users.UUID = pending_requests.UUID").Scan(&records)
	return records

}

func GetDataofPendingRequestsByUUID(DB *gorm.DB, uuid int) []types.PendingRequestData {
	var records []types.PendingRequestData
	DB.Table("pending_requests").Select("users.UUID, books.B_UID as BUID, users.NAME as USER_NAME, books.NAME as BOOK_NAME, pending_requests.TYPE").Joins("join books on books.B_UID = pending_requests.B_UID").Joins("join users on users.UUID = pending_requests.UUID").Where(&PENDING_REQUESTS{UUID: uuid}).Scan(&records)
	return records
}

func DeletePendingRequest(DB *gorm.DB, buid int, uuid int) {
	var req PENDING_REQUESTS
	DB.Where(&PENDING_REQUESTS{BUID: buid, UUID: uuid}).Delete(&req)
}

func UpdateUserData(DB *gorm.DB, uuid int, name string, phn_no int64) {
	var user USER
	res := DB.Where("UUID=?", uuid).First(&user)
	utils.CheckNilErr(res.Error, "Error with DB while finding user") //// if an unknown error occurs here then panic will stop the control flow, if checkNilErr has to be removed then adequate return should be added
	user.NAME = name
	user.PHN_NO = phn_no
	DB.Save(&user)

}

func DeleteBook(DB *gorm.DB, buid int) {
	var book BOOKS
	DB.Where(&BOOKS{BUID: buid}).Delete(&book)
}

func CreateNewBook(DB *gorm.DB, name string, desc string, qty int) {
	book := BOOKS{
		NAME:        name,
		DESCRIPTION: desc,
		CHECKIN:     qty,
		TOTAL:       qty,
	}
	res := DB.Create(&book)
	utils.CheckNilErr(res.Error, "Unable to create book")
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
	utils.CheckNilErr(res.Error, "Error with DB while finding user") //// if an unknown error occurs here then panic will stop the control flow, if checkNilErr has to be removed then adequate return should be added
	user.ADMIN_REQUEST = &truevalue
	DB.Save(&user)

}

func SetAdminReq(DB *gorm.DB, uuid int, status bool) {
	var user USER
	falsevalue := false
	res := DB.Where("UUID=?", uuid).First(&user)
	utils.CheckNilErr(res.Error, "Error with DB while finding user") //// if an unknown error occurs here then panic will stop the control flow, if checkNilErr has to be removed then adequate return should be added
	if status {
		user.ADMIN_REQUEST = nil
		user.ROLE = "admin"
	} else {
		user.ADMIN_REQUEST = &falsevalue
	}
	DB.Save(&user)
}
