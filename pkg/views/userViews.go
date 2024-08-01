package views

import "html/template"

func UserDashView() *template.Template {

	files := []string{
		"pkg/templates/userDash.html",

		"pkg/templates/partials/userTopbar.html",
		"pkg/templates/partials/userSidebar.html",
	}
	temp := template.Must(template.New("userDash.html").Funcs(funcMap).ParseFiles(files...))
	return temp
}

func UserCvtAdminView() *template.Template {

	files := []string{
		"pkg/templates/userCvtAdmin.html",

		"pkg/templates/partials/userTopbar.html",
		"pkg/templates/partials/userSidebar.html",
	}
	temp := template.Must(template.New("userCvtAdmin.html").Funcs(funcMap).ParseFiles(files...))
	return temp
}
