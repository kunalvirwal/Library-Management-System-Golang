package models

import (
	"github.com/kunalvirwal/go-mvc/pkg/types"
	"github.com/kunalvirwal/go-mvc/pkg/utils"
	"gorm.io/gorm"
)

func GetAllPendingCheckinReqByUUID(DB *gorm.DB, uuid int) []PENDING_REQUESTS {
	var reqs []PENDING_REQUESTS
	DB.Where(&PENDING_REQUESTS{UUID: uuid, TYPE: true}).Find(&reqs)
	return reqs
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
			utils.CheckNilErr(res.Error, "Error with DB while finding pending request")
		}
	}
	return req, true
}

func CreateNewBookReq(DB *gorm.DB, buid int, uuid int, typ bool) {
	req := PENDING_REQUESTS{UUID: uuid, BUID: buid, TYPE: typ}
	res := DB.Create(req)
	if res.Error != nil {
		utils.CheckNilErr(res.Error, "Unable to create book Pending request ")
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
