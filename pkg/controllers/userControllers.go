package controllers

import (
	"net/http"
	"slices"

	"github.com/kunalvirwal/go-mvc/pkg/models"
	"github.com/kunalvirwal/go-mvc/pkg/types"
	"github.com/kunalvirwal/go-mvc/pkg/utils"
	"github.com/kunalvirwal/go-mvc/pkg/views"
)

func GetUserData(w http.ResponseWriter, r *http.Request) {
	data, ok := r.Context().Value(user).(*types.Claims)

	if data != nil && ok {
		db, err := models.Connection()
		utils.CheckNilErr(err, "Unable To load DB in GetUserData")
		books := models.GetCheckedOutBooks(db, data.UUID)
		result := models.GetAllPendingCheckinReqByUUID(db, data.UUID)
		past_books := models.GetAllPastCheckedInBooks(db, data.UUID)
		slices.Reverse(past_books)

		for i := 0; i < len(books); i++ {
			for j := 0; j < len(result); j++ {
				if result[j].BUID == books[i].BUID {
					books[i].Req = true
					break
				}
			}
		}

		sendData := map[string]interface{}{
			"name":       data.Name,
			"books":      books,
			"past_books": past_books,
			"path":       r.URL.String(),
		}

		t := views.UserDashView()
		t.Execute(w, sendData)
	} else {
		http.Redirect(w, r, "/", http.StatusBadRequest)
	}
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	data, ok := r.Context().Value(user).(*types.Claims)

	if data != nil && ok {
		uuid := data.UUID
		db, err := models.Connection()
		utils.CheckNilErr(err, "Unable To load DB in GetUserData")
		user, found := models.SearchUserUUID(db, uuid)

		if !found {
			http.Redirect(w, r, "/", http.StatusBadRequest)
			return
		}

		if user.ROLE == utils.Admin {
			http.Redirect(w, r, "/logout", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/cvt_admin", http.StatusSeeOther)
		return
	} else {
		http.Redirect(w, r, "/", http.StatusBadRequest)
	}
}
