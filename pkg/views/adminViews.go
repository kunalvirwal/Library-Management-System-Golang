package views

import "html/template"

func AdminDashView() *template.Template {
	files := []string{
		"pkg/templates/adminDash.html",
		// "pkg/templates/basic.html",

		"pkg/templates/partials/adminTopbar.html",
		"pkg/templates/partials/adminSidebar.html",
	}

	temp := template.Must(template.New("adminDash.html").Funcs(funcMap).ParseFiles(files...))
	return temp
}

func EditbookView() *template.Template {
	files := []string{
		"pkg/templates/editBookPage.html",
		// "pkg/templates/basic.html",

		"pkg/templates/partials/adminTopbar.html",
		"pkg/templates/partials/adminSidebar.html",
	}

	temp := template.Must(template.New("editBookPage.html").Funcs(funcMap).ParseFiles(files...))
	return temp
}

func NewBookView() *template.Template {
	files := []string{
		"pkg/templates/CreateBook.html",
		// "pkg/templates/basic.html",

		"pkg/templates/partials/adminTopbar.html",
		"pkg/templates/partials/adminSidebar.html",
	}

	temp := template.Must(template.New("CreateBook.html").Funcs(funcMap).ParseFiles(files...))
	return temp
}
