package views

import "html/template"

func LoginView() *template.Template {
	temp := template.Must(template.ParseFiles("pkg/templates/login.html"))
	return temp
}

func SignupView() *template.Template {
	temp := template.Must(template.ParseFiles("pkg/templates/signup.html"))
	return temp
}
