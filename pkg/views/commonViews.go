package views

import (
	"html/template"
	"time"
)

var funcMap = template.FuncMap{
	"add": func(a int, b int) int {
		return a + b
	},
	"sub": func(a int, b int) int {
		return a - b
	},
	"ptr": func(t *time.Time) string {
		if t != nil {
			return t.Format("2006-01-02 15:04:05") // Format the date as needed
		}
		return "N/A"
	},
	"date": func(t time.Time) string {
		return t.Format("2006-01-02 15:04:05") // Format the date as needed
	},
}

func LoginView() *template.Template {
	temp := template.Must(template.ParseFiles("pkg/templates/login.html"))
	return temp
}

func SignupView() *template.Template {
	temp := template.Must(template.ParseFiles("pkg/templates/signup.html"))
	return temp
}

func BookCatalogView() *template.Template {

	files := []string{
		"pkg/templates/bookCatalog.html",
		// "pkg/templates/basic.html",

		"pkg/templates/partials/userTopbar.html",
		"pkg/templates/partials/userSidebar.html",
		"pkg/templates/partials/adminTopbar.html",
		"pkg/templates/partials/adminSidebar.html",
	}

	temp := template.Must(template.New("bookCatalog.html").Funcs(funcMap).ParseFiles(files...))
	return temp
}

func BookPageView() *template.Template {

	files := []string{
		"pkg/templates/bookPage.html",
		// "pkg/templates/basic.html",

		"pkg/templates/partials/userTopbar.html",
		"pkg/templates/partials/userSidebar.html",
		"pkg/templates/partials/adminTopbar.html",
		"pkg/templates/partials/adminSidebar.html",
	}

	temp := template.Must(template.New("bookPage.html").Funcs(funcMap).ParseFiles(files...))
	return temp
}

func PendingView() *template.Template {

	files := []string{
		"pkg/templates/pendingReq.html",
		// "pkg/templates/basic.html",

		"pkg/templates/partials/userTopbar.html",
		"pkg/templates/partials/userSidebar.html",
		"pkg/templates/partials/adminTopbar.html",
		"pkg/templates/partials/adminSidebar.html",
	}

	temp := template.Must(template.New("pendingReq.html").Funcs(funcMap).ParseFiles(files...))
	return temp
}

func AccountView() *template.Template {

	files := []string{
		"pkg/templates/account.html",
		// "pkg/templates/basic.html",

		"pkg/templates/partials/userTopbar.html",
		"pkg/templates/partials/userSidebar.html",
		"pkg/templates/partials/adminTopbar.html",
		"pkg/templates/partials/adminSidebar.html",
	}

	temp := template.Must(template.New("account.html").Funcs(funcMap).ParseFiles(files...))
	return temp
}
