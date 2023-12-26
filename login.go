package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func HasAuthCookie(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if isLoggedIn(w,r) {
			next(w, r, ps)
		} else {
			fmt.Println("nah")
			tmpl.ExecuteTemplate(w, "nonloginhome.html", nil)
		}
	}
}

func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "exampleCookie",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}

func registerHandler(response http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	switch request.Method {
	case "GET" : 
		tmpl.ExecuteTemplate(response,"register.html",nil)
	case "POST" : 
	username := request.FormValue("username")
	password := request.FormValue("password")
	fmt.Println(username, "  " ,password)
	if username != "" && password != "" {
		rdxHset("user",username,password)
		val, _ := rdxHget("user",username)
		fmt.Println(val)
		http.Redirect(response, request, "/login", 302)
	} else {
		http.Redirect(response, request, "/register", 302)
	}
	default:
		fmt.Fprintf(response,"Method Not Allowed")
	}
}

// login handler

// index page
const already = `<!DOCTYPE HTML><p>Already Logged In</p>`

func loginHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	switch r.Method {
	case "GET" : {
		if isLoggedIn(w,r) {
			fmt.Fprintf(w, already)
		} else {
			tmpl.ExecuteTemplate(w,"login.html",nil)
		}
	}
	case "POST" : {
	username := r.FormValue("username")
	pass := r.FormValue("password")
	if username != "" && pass != "" {
		val, _ := rdxHget("user",username)
		fmt.Println(val)
		fmt.Println(pass)
		if val == pass {
			setSession(username,w)
			http.Redirect(w, r, "/", 301)
		} else {
			fmt.Fprintf(w,"Wrong credentials")
		}
	}
	}
	default:
		fmt.Fprintf(w,"Method Not Allowed")
	}
}

// logout handler

func logoutHandler(response http.ResponseWriter, request *http.Request, _ httprouter.Params) {
		cookie := &http.Cookie{
		Name:   "exampleCookie",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
	http.Redirect(response, request, "/", 302)
}


func indexPageHandler(response http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	tmpl.ExecuteTemplate(response,"login.html",nil)
}

func setSession(name string, w http.ResponseWriter) {
    usernameCookie := http.Cookie{
        Name:     "exampleCookie",
        Value:    name,
        Path:     "/",
        MaxAge:   3600,
        HttpOnly: true,
        Secure:   true,
        SameSite: http.SameSiteLaxMode,
    }
    http.SetCookie(w, &usernameCookie)
}

func isLoggedIn(w http.ResponseWriter, r *http.Request) bool {
    cookie, err := r.Cookie("exampleCookie")
    if err != nil {
        return false
    }
	//~ setSession(cookie.Value,w)
	fmt.Println(cookie.Value)
	return true
}