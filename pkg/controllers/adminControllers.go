package controllers

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/kunalvirwal/go-mvc/pkg/models"
	"github.com/kunalvirwal/go-mvc/pkg/types"
	"github.com/kunalvirwal/go-mvc/pkg/utils"
	"github.com/kunalvirwal/go-mvc/pkg/views"
)

func GetAdminData(w http.ResponseWriter, r *http.Request) {
	data, ok := r.Context().Value(user).(*types.Claims)
	if data != nil && ok {
		db, err := models.Connection()
		utils.CheckNilErr(err, "Unable To load DB in GetUserData")
		////////////////chek if a user with this jwt exists
		dataset := models.GetAllBooks(db)
		checked_out := 0
		for _, val := range dataset {
			if val.CHECKIN == 0 {
				checked_out++
			}
		}
		users := models.GetAllUsers(db)
		pending_requests := models.GetAllPendingRequests(db)
		books := models.GetCheckedOutBooks(db, data.UUID)
		past_books := models.GetAllPastCheckedInBooks(db, data.UUID)
		slices.Reverse(past_books)
		overdue := 1 // use borrowing history for implementing overdues
		sendData := map[string]interface{}{
			"name":             data.Name,
			"no_of_books":      len(dataset),
			"checked_out":      checked_out,
			"no_of_users":      len(users),
			"pending_requests": len(pending_requests),
			"books":            books,
			"past_books":       past_books,
			"overdue":          overdue,
			"path":             r.URL.String(),
		}

		t := views.AdminDashView()
		t.Execute(w, sendData)
	} else {
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
	}

}

func Approve(w http.ResponseWriter, r *http.Request) {
	data, ok := r.Context().Value(user).(*types.Claims)
	if data != nil && ok {
		params := mux.Vars(r)
		inp_uuid, err := strconv.ParseInt(strings.TrimSpace(params["uuid"]), 10, 64)
		inp_buid, err2 := strconv.ParseInt(strings.TrimSpace(params["buid"]), 10, 64)
		if err != nil || inp_uuid <= 0 || err2 != nil || inp_buid <= 0 {
			http.Redirect(w, r, "/pending", http.StatusSeeOther)
			return
		} else {

			db, err := models.Connection()
			utils.CheckNilErr(err, "Unable To load DB in GetUserData")
			book, found := models.GetBook(db, int(inp_buid))
			req, req_exist := models.PendingReqExist(db, int(inp_buid), int(inp_uuid))
			if !found || !req_exist {
				http.Redirect(w, r, "/pending", http.StatusSeeOther)
				return
			}

			if req.TYPE { // checkin req
				models.DeletePendingRequest(db, req.BUID, req.UUID)
				if book.CHECKIN < book.TOTAL {
					models.UpdateCheckinEntry(db, int(inp_buid), int(inp_uuid))
					models.EditBook(db, int(inp_buid), book.NAME, book.DESCRIPTION, (book.CHECKIN + 1), book.TOTAL)
				} else {
					denyURL := "/admin/deny/" + strconv.Itoa(int(inp_uuid)) + "/" + strconv.Itoa(int(inp_buid))
					http.Redirect(w, r, denyURL, http.StatusSeeOther)
				}

			} else { // checkout req
				models.DeletePendingRequest(db, req.BUID, req.UUID)
				if book.CHECKIN > 0 {
					models.CreateCheckoutEntry(db, int(inp_buid), int(inp_uuid))
					models.EditBook(db, int(inp_buid), book.NAME, book.DESCRIPTION, (book.CHECKIN - 1), book.TOTAL)

				} else {
					denyURL := "/admin/deny/" + strconv.Itoa(int(inp_uuid)) + "/" + strconv.Itoa(int(inp_buid))
					http.Redirect(w, r, denyURL, http.StatusSeeOther)
				}

			}
			if data.Role == "admin" {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return ///////////////////////////////////////////////////////// redirct can happen to / if in bookpage and to /pending in pending requests
			}

		}
	} else {
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
	}
}

func Deny(w http.ResponseWriter, r *http.Request) {
	data, ok := r.Context().Value(user).(*types.Claims)
	if data != nil && ok {
		params := mux.Vars(r)
		inp_uuid, err := strconv.ParseInt(strings.TrimSpace(params["uuid"]), 10, 64)
		inp_buid, err2 := strconv.ParseInt(strings.TrimSpace(params["buid"]), 10, 64)

		if err != nil || inp_uuid <= 0 || err2 != nil || inp_buid <= 0 {
			http.Redirect(w, r, "/pending", http.StatusSeeOther)
			return
		} else {
			db, err := models.Connection()
			utils.CheckNilErr(err, "Unable To load DB in GetUserData")
			_, found := models.GetBook(db, int(inp_buid))
			req, req_exist := models.PendingReqExist(db, int(inp_buid), int(inp_uuid))
			if !found || !req_exist {
				http.Redirect(w, r, "/pending", http.StatusSeeOther)
				return
			}
			models.DeletePendingRequest(db, req.BUID, req.UUID)

			http.Redirect(w, r, "/pending", http.StatusSeeOther)

		}
	} else {
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
	}
}

func Editbook(w http.ResponseWriter, r *http.Request) {
	data, ok := r.Context().Value(user).(*types.Claims)
	if data != nil && ok {
		params := mux.Vars(r)
		inp_buid, err := strconv.ParseInt(strings.TrimSpace(params["buid"]), 10, 64)

		if err != nil || inp_buid <= 0 {
			http.Redirect(w, r, "/books", http.StatusSeeOther)
			return
		} else {

			db, err := models.Connection()
			utils.CheckNilErr(err, "Unable To load DB in GetUserData")
			book, found := models.GetBook(db, int(inp_buid))
			if !found {
				http.Redirect(w, r, "/books", http.StatusSeeOther)
				return
			}
			sendData := map[string]interface{}{
				"book": book,
				"path": r.URL.String(),
			}

			t := views.EditbookView()
			t.Execute(w, sendData)

		}
	} else {
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
	}
}

func SaveEditbookChanges(w http.ResponseWriter, r *http.Request) {
	data, ok := r.Context().Value(user).(*types.Claims)
	if data != nil && ok {
		params := mux.Vars(r)
		inp_buid, err := strconv.ParseInt(strings.TrimSpace(params["buid"]), 10, 64)

		if err != nil || inp_buid <= 0 {
			http.Redirect(w, r, "/books", http.StatusSeeOther)
			return
		} else {
			err := r.ParseForm()
			if err != nil {
				http.Redirect(w, r, "/logout", http.StatusSeeOther)
				fmt.Println("Invalid Post req paramemeters")
				return
			}
			db, err := models.Connection()
			utils.CheckNilErr(err, "Unable To load DB in GetUserData")
			book, found := models.GetBook(db, int(inp_buid))
			if !found {
				http.Redirect(w, r, "/books", http.StatusSeeOther)
				return
			}
			name := strings.TrimSpace(r.FormValue("name"))
			desc := strings.TrimSpace(r.FormValue("description"))

			qty, err := strconv.ParseInt(strings.TrimSpace(r.FormValue("qty")), 10, 64)

			checked_out := book.TOTAL - book.CHECKIN
			if name == "" || len(name) > 50 || !utils.Sanitise(name, false) || desc == "" || len(desc) > 2000 || err != nil || int(qty) < checked_out {
				bookURL := "/books/" + strconv.Itoa(book.BUID)
				http.Redirect(w, r, bookURL, http.StatusSeeOther)
				return
			}
			decrease := book.TOTAL - int(qty)
			decreased_checkin := book.CHECKIN - decrease //decrease can only happen from books checked in
			models.EditBook(db, book.BUID, name, desc, decreased_checkin, int(qty))

			bookURL := "/books/" + strconv.Itoa(book.BUID)
			http.Redirect(w, r, bookURL, http.StatusSeeOther)

		}
	} else {
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
	}
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	data, ok := r.Context().Value(user).(*types.Claims)
	if data != nil && ok {
		params := mux.Vars(r)
		inp_buid, err := strconv.ParseInt(strings.TrimSpace(params["buid"]), 10, 64)

		if err != nil || inp_buid <= 0 {
			http.Redirect(w, r, "/books", http.StatusSeeOther)
			return
		} else {

			db, err := models.Connection()
			utils.CheckNilErr(err, "Unable To load DB in GetUserData")
			book, found := models.GetBook(db, int(inp_buid))
			if !found {
				http.Redirect(w, r, "/books", http.StatusSeeOther)
				return
			}
			if book.CHECKIN < book.TOTAL {
				bookURL := "/books/" + strconv.Itoa(book.BUID)
				http.Redirect(w, r, bookURL, http.StatusSeeOther)
				return
			}
			models.DeleteBook(db, book.BUID)
			http.Redirect(w, r, "/books", http.StatusSeeOther)

		}
	} else {
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
	}
}

func NewBook(w http.ResponseWriter, r *http.Request) {
	data, ok := r.Context().Value(user).(*types.Claims)
	if data != nil && ok {

		t := views.NewBookView()
		t.Execute(w, map[string]interface{}{"path": r.URL.String()})

	} else {
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
	}
}

func SaveNewBook(w http.ResponseWriter, r *http.Request) {
	data, ok := r.Context().Value(user).(*types.Claims)
	if data != nil && ok {

		err := r.ParseForm()
		if err != nil {
			http.Redirect(w, r, "/logout", http.StatusSeeOther)
			fmt.Println("Invalid Post req paramemeters")
			return
		}
		name := strings.TrimSpace(r.FormValue("name"))
		desc := strings.TrimSpace(r.FormValue("description"))
		qty, err := strconv.ParseInt(r.FormValue("qty"), 10, 64)
		if err != nil || !utils.Sanitise(name, false) || qty < 1 || len(name) == 0 || len(name) > 50 || len(desc) > 2000 || len(desc) == 0 {
			http.Redirect(w, r, "/admin/addbook", http.StatusSeeOther)
			return
		}
		db, err := models.Connection()
		utils.CheckNilErr(err, "Unable to create Db instance")
		models.CreateNewBook(db, name, desc, int(qty))
		http.Redirect(w, r, "/books", http.StatusSeeOther)

	} else {
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
	}
}
