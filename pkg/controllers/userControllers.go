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
		////////////////chek if a user with this jwt exists????
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
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
	}
}
