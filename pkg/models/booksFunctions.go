package models

import (
	"github.com/kunalvirwal/go-mvc/pkg/types"
	"github.com/kunalvirwal/go-mvc/pkg/utils"
	"gorm.io/gorm"
)

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

func GetAllBooks(DB *gorm.DB) []BOOKS {
	var dataset []BOOKS
	DB.Find(&dataset)
	return dataset
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
			utils.CheckNilErr(res.Error, "Error with DB while finding book") // if an unknown error occurs here then panic will stop the control flow, if checkNilErr has to be removed then adequate return should be added
		}
	}
	return book, true
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
