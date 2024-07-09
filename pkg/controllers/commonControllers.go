package controllers

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/kunalvirwal/go-mvc/pkg/models"
	"github.com/kunalvirwal/go-mvc/pkg/types"
	"github.com/kunalvirwal/go-mvc/pkg/utils"
	"github.com/kunalvirwal/go-mvc/pkg/views"
)

var user types.ContextKeyType = "abc"

func LoginPage(tried bool) http.HandlerFunc { ////////////////also check for pre
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, ok := r.Context().Value(user).(*types.Claims)
		if ok && data != nil {
			if data.Role == "admin" {
				http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
				return
			}
			if data.Role == "user" {
				http.Redirect(w, r, "/user/dashboard", http.StatusSeeOther)
				return
			}
		}
		// fmt.Println("rendering login page")
		t := views.LoginView()
		t.Execute(w, tried)
	})
}

func SignupPage(exists bool) http.HandlerFunc { ////////////////also check for pre
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		t := views.SignupView()
		t.Execute(w, exists)
	})
}

func NewUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
		fmt.Println("Invalid Post req paramemeters")
		return
	}
	name := strings.TrimSpace(r.FormValue("name"))
	password := strings.TrimSpace(r.FormValue("password"))
	email := strings.TrimSpace(r.FormValue("email"))
	phn_no, err := strconv.ParseInt(r.FormValue("phn_no"), 10, 64)
	if err != nil || !utils.Sanitise(name, false) || len([]byte(password)) > 72 || phn_no > 9999999999 || phn_no < 1000000000 || len(name) == 0 || len(name) > 50 {
		t := views.SignupView()
		t.Execute(w, true)
		return
	}
	password = utils.SaltNhash(password)

	db, err := models.Connection()
	utils.CheckNilErr(err, "Unable to create Db instance")
	_, found := models.SearchUserEmail(db, email)
	if found {
		t := views.SignupView()
		t.Execute(w, true)
		return
	}
	models.CreateNewUser(db, name, email, phn_no, password)
	Logging(w, r)
}

func Logging(w http.ResponseWriter, r *http.Request) {
	password := strings.TrimSpace(r.FormValue("password"))
	email := strings.TrimSpace(r.FormValue("email"))
	db, err := models.Connection()
	utils.CheckNilErr(err, "Unable to create Db instance")
	user, found := models.SearchUserEmail(db, email)
	correctPwd := utils.MatchHashtoPassword(user.PASSWORD, password)
	if found && correctPwd {
		// fmt.Println(user.UUID)
		token, expirationTime := utils.GenerateJWT(user.UUID, email, user.NAME, user.ROLE)

		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   token,
			Expires: expirationTime,
		})

		fmt.Println("JWT set")
		http.Redirect(w, r, "/", http.StatusSeeOther)

	} else {
		t := views.LoginView()
		t.Execute(w, true)
	}
}

func GetBookCatalog(w http.ResponseWriter, r *http.Request) {
	data1, ok := r.Context().Value(user).(*types.Claims)
	if data1 != nil && ok {
		db, err := models.Connection()
		utils.CheckNilErr(err, "Unable to create Db instance")

		dataset := models.GetAllBooks(db)
		var datasetNew []models.BOOKS
		role := data1.Role
		books_per_page := 5

		query := r.URL.Query()
		var search string

		if query["search"] != nil {
			search := string(query["search"][0])
			// fmt.Println(search, final)
			search = strings.TrimSpace(search)
			for _, val := range dataset {

				if strings.Contains(val.NAME, search) {
					datasetNew = append(datasetNew, val)
				}
			}
			dataset = datasetNew
		}

		no_of_pages := int(math.Ceil(float64(len(dataset)) / float64(books_per_page)))
		var final int
		if query["page"] != nil {
			page, err := strconv.ParseInt(query["page"][0], 10, 64)
			if err != nil || page <= 0 || page > int64(no_of_pages) {
				final = 1
			} else {
				final = int(page)
			}
		} else {
			final = 1
		}

		start_index := (final - 1) * books_per_page
		end_index := final * books_per_page

		if end_index > len(dataset) {
			end_index = len(dataset)
		}
		sendData := map[string]interface{}{
			"data":           dataset[start_index:end_index],
			"page":           final,
			"books_per_page": books_per_page,
			"no_of_pages":    make([]int, no_of_pages),
			"start_index":    start_index,
			"end_index":      end_index,
			"search":         search,
			"role":           role,
			"path":           r.URL.String(),
		}

		t := views.BookCatalogView()
		t.Execute(w, sendData)
	} else {
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
	}

}

func GetBookPage(w http.ResponseWriter, r *http.Request) {
	data, ok := r.Context().Value(user).(*types.Claims)
	if data != nil && ok {
		params := mux.Vars(r)
		inp_buid, err := strconv.ParseInt(strings.TrimSpace(params["buid"]), 10, 64)
		if err != nil || inp_buid <= 0 {
			http.Redirect(w, r, "/books", http.StatusSeeOther)
			return
		} else {
			db, err := models.Connection()
			utils.CheckNilErr(err, "Unable to create Db instance")
			book, found := models.GetBook(db, int(inp_buid))
			if !found {
				http.Redirect(w, r, "/books", http.StatusSeeOther)
				return
			}
			// fmt.Println(book)
			_, found2 := models.PendingReqExist(db, int(inp_buid), data.UUID)
			user_has := (models.IsCheckedOutByUser(db, int(inp_buid), data.UUID) || found2)
			role := data.Role

			sendData := map[string]interface{}{
				"book":     book,
				"role":     role,
				"user_has": user_has,
				"path":     r.URL.String(),
			}
			t := views.BookPageView()
			t.Execute(w, sendData)
		}
	} else {
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
	}

}

func MakeCheckoutReq(w http.ResponseWriter, r *http.Request) {
	data, ok := r.Context().Value(user).(*types.Claims)
	err := r.ParseForm()
	if data != nil && ok && err == nil {

		inp_val := strings.TrimSpace(r.FormValue("buid"))
		inp_buid, err := strconv.ParseInt(inp_val, 10, 64)
		if err != nil || inp_buid <= 0 {
			http.Redirect(w, r, "/books", http.StatusSeeOther)
			return
		}
		db, err := models.Connection()
		utils.CheckNilErr(err, "Unable to create Db instance")

		user_has := models.IsCheckedOutByUser(db, int(inp_buid), data.UUID)
		book, found := models.GetBook(db, int(inp_buid))
		_, req_exist := models.PendingReqExist(db, int(inp_buid), data.UUID)

		if user_has || !found || req_exist || book.CHECKIN == 0 {
			http.Redirect(w, r, "/books", http.StatusSeeOther)
			return
		}

		models.CreateNewBookReq(db, int(inp_buid), data.UUID, false)

		if data.Role == "admin" {

			approveURL := "/admin/approve/" + strconv.Itoa(data.UUID) + "/" + strconv.Itoa(int(inp_buid))
			http.Redirect(w, r, approveURL, http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/pending", http.StatusSeeOther)

	} else {
		fmt.Println("Invalid Post req paramemeters")
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
	}
}

func MakeCheckinReq(w http.ResponseWriter, r *http.Request) {
	data, ok := r.Context().Value(user).(*types.Claims)
	err := r.ParseForm()
	if data != nil && ok && err == nil {

		inp_val := strings.TrimSpace(r.FormValue("buid"))
		inp_buid, err := strconv.ParseInt(inp_val, 10, 64)
		if err != nil || inp_buid <= 0 {
			http.Redirect(w, r, "/books", http.StatusSeeOther)
			return
		}
		db, err := models.Connection()
		utils.CheckNilErr(err, "Unable to create Db instance")

		user_has := models.IsCheckedOutByUser(db, int(inp_buid), data.UUID)
		book, found := models.GetBook(db, int(inp_buid))
		_, req_exist := models.PendingReqExist(db, int(inp_buid), data.UUID)

		if !user_has || !found || req_exist || book.CHECKIN == book.TOTAL {
			http.Redirect(w, r, "/", http.StatusSeeOther)

			return
		}

		models.CreateNewBookReq(db, int(inp_buid), data.UUID, true)

		if data.Role == "admin" {
			approveURL := "/admin/approve/" + strconv.Itoa(data.UUID) + "/" + strconv.Itoa(int(inp_buid))
			http.Redirect(w, r, approveURL, http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/pending", http.StatusSeeOther)

	} else {
		fmt.Println("Invalid Post req paramemeters")
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
	}
}

func GetPending(w http.ResponseWriter, r *http.Request) {
	data, ok := r.Context().Value(user).(*types.Claims)
	if data != nil && ok {

		role := data.Role
		uuid := data.UUID
		db, err := models.Connection()
		utils.CheckNilErr(err, "Unable to create Db instance")
		var records, checkins, checkouts []types.PendingRequestData

		if role == "admin" {
			records = models.GetDataofAllPendingRequests(db)
		} else {
			records = models.GetDataofPendingRequestsByUUID(db, uuid)
		}
		// fmt.Println(records)
		for _, val := range records {
			if val.TYPE {
				checkins = append(checkins, val)
			} else {
				checkouts = append(checkouts, val)
			}
		}

		sendData := map[string]interface{}{
			"role":      role,
			"checkins":  checkins,
			"checkouts": checkouts,
			"path":      r.URL.String(),
		}

		t := views.PendingView()
		t.Execute(w, sendData)

	} else {
		fmt.Println("Invalid Post req paramemeters")
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
	}
}

func Account(w http.ResponseWriter, r *http.Request) {
	data, ok := r.Context().Value(user).(*types.Claims)
	if data != nil && ok {

		Email := data.Email
		db, err := models.Connection()
		utils.CheckNilErr(err, "Unable to create Db instance")
		user, found := models.SearchUserEmail(db, Email)
		if !found {
			http.Redirect(w, r, "/logout", http.StatusSeeOther)
			return
		}
		sendData := map[string]interface{}{
			"UUID":   user.UUID,
			"NAME":   user.NAME,
			"EMAIL":  user.EMAIL,
			"PHN_NO": user.PHN_NO,
			"role":   user.ROLE,
			"path":   r.URL.String(),
		}

		t := views.AccountView()
		t.Execute(w, sendData)

	} else {
		fmt.Println("Invalid Post req paramemeters")
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
	}

}

func EditAccount(w http.ResponseWriter, r *http.Request) {
	data, ok := r.Context().Value(user).(*types.Claims)
	if data != nil && ok {

		err := r.ParseForm()
		if err != nil {
			http.Redirect(w, r, "/logout", http.StatusSeeOther)
			fmt.Println("Invalid Post req paramemeters")
			return
		}
		name := strings.TrimSpace(r.FormValue("name"))
		email := data.Email
		uuid := data.UUID
		phn_no, err := strconv.ParseInt(strings.TrimSpace(r.FormValue("phn_no")), 10, 64)
		if err != nil || !utils.Sanitise(name, false) || phn_no > 9999999999 || phn_no < 1000000000 || len(name) == 0 || len(name) > 50 {
			http.Redirect(w, r, "/accounts", http.StatusSeeOther)
			return
		}
		db, err := models.Connection()
		utils.CheckNilErr(err, "Unable to create Db instance")
		_, found := models.SearchUserEmail(db, email)
		if !found {
			http.Redirect(w, r, "/logout", http.StatusSeeOther)
			return
		}

		models.UpdateUserData(db, uuid, name, phn_no)
		http.Redirect(w, r, "/logout", http.StatusSeeOther)

	} else {
		fmt.Println("Invalid Post req paramemeters")
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
	}
}

func GetCvtAdmin(w http.ResponseWriter, r *http.Request) {
	data, ok := r.Context().Value(user).(*types.Claims)
	if data != nil && ok {
		role := data.Role
		email := data.Email
		uuid := data.UUID

		db, err := models.Connection()
		utils.CheckNilErr(err, "Unable to create Db instance")

		user, found := models.SearchUserEmail(db, email)
		if !found {
			http.Redirect(w, r, "/logout", http.StatusSeeOther)
		}

		if role == "user" {
			applied := user.ADMIN_REQUEST
			if applied == nil {

				Refresh(w, r)
				return
			}

			sendData := map[string]interface{}{
				"applied": *applied,
				"uuid":    uuid,
				"path":    r.URL.String(),
			}
			t := views.UserCvtAdminView()
			t.Execute(w, sendData)

		} else if role == "admin" {
			users := models.GetAdminRequests(db)
			sendData := map[string]interface{}{
				"users": users,
				"path":  r.URL.String(),
			}

			t := views.AdminCvtAdminView()
			t.Execute(w, sendData)

		}

	} else {
		fmt.Println("Invalid Post req paramemeters")
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
	}
}

func PostCvtAdmin(w http.ResponseWriter, r *http.Request) {
	data, ok := r.Context().Value(user).(*types.Claims)
	if data != nil && ok {
		role := data.Role
		email := data.Email
		err := r.ParseForm()
		if err != nil {
			http.Redirect(w, r, "/logout", http.StatusSeeOther)
			fmt.Println("Invalid Post req paramemeters")
			return
		}

		db, err := models.Connection()
		utils.CheckNilErr(err, "Unable to create Db instance")

		ourUser, found := models.SearchUserEmail(db, email)
		if !found {
			http.Redirect(w, r, "/logout", http.StatusSeeOther)
		}

		if role == "user" {
			if ourUser.ADMIN_REQUEST == nil || *ourUser.ADMIN_REQUEST {
				http.Redirect(w, r, "/cvt_admin", http.StatusSeeOther)
			} else {
				models.CreateAdminReq(db, ourUser.UUID)
				http.Redirect(w, r, "/cvt_admin", http.StatusSeeOther)
			}
		} else if role == "admin" {
			err := r.ParseForm()
			if err != nil {
				http.Redirect(w, r, "/logout", http.StatusSeeOther)
				fmt.Println("Invalid Post req paramemeters")
				return
			}
			inp_uuid := strings.TrimSpace(r.FormValue("approve"))
			status := true
			if inp_uuid == "" {
				inp_uuid = strings.TrimSpace(r.FormValue("deny"))
				status = false
			}
			target_uuid, err := strconv.ParseInt(inp_uuid, 10, 64)
			if err != nil || target_uuid <= 0 {
				http.Redirect(w, r, "/cvt_admin", http.StatusSeeOther)
				return
			}
			user, found := models.SearchUserUUID(db, int(target_uuid))
			if !found {
				fmt.Println("Invalid Post req paramemeters")
				http.Redirect(w, r, "/logout", http.StatusSeeOther)
			}
			if *user.ADMIN_REQUEST {
				models.SetAdminReq(db, user.UUID, status)

			}
			http.Redirect(w, r, "/cvt_admin", http.StatusSeeOther)
			return

		}

	} else {
		fmt.Println("Invalid Post req paramemeters")
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	utils.DeleteJWT(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
