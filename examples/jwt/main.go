package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"log"

	"github.com/golang-jwt/jwt/v5"
)

type jwtCustomClaims struct {
	jwt.RegisteredClaims
}

var key = []byte("secret")

func main() {

	s := http.NewServeMux()
	s.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
	})
	s.HandleFunc("GET /login", handleShowLogin)
	s.HandleFunc("POST /login", handlePostLogin)
	http.ListenAndServe(":8083", s)
}

func handleShowLogin(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("err") != "" {
		tmpl := template.Must(template.ParseFiles("login.html"))
		err := tmpl.Execute(w, map[string]any{
			"error": r.URL.Query().Get("err"),
		})
		if err != nil {
			w.Write([]byte(fmt.Sprintf("error: %s\n", err)))
			return
		}
		return
	}

	tmpl := template.Must(template.ParseFiles("login.html"))
	err := tmpl.Execute(w, map[string]any{
		"error": "",
	})
	if err != nil {
		w.Write([]byte(fmt.Sprintf("error: %s\n", err)))
		return
	}
}

func handlePostLogin(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("username") != "test" || r.FormValue("password") != "test" {
		tmpl := template.Must(template.ParseFiles("login.html"))
		err := tmpl.Execute(w, map[string]any{
			"error": "Invalid username or password",
		})
		if err != nil {
			w.Write([]byte(fmt.Sprintf("error: %s\n", err)))
			return
		}
		return
	}

	claims := &jwtCustomClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 1)),
			Issuer:    "haproxy",
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	s, err := t.SignedString(key)
	if err != nil {
		tmpl := template.Must(template.ParseFiles("login.html"))
		err := tmpl.Execute(w, map[string]any{
			"error": err.Error(),
		})
		if err != nil {
			w.Write([]byte(fmt.Sprintf("error: %s\n", err)))
			return
		}
		return
	}
	// set cookie that expires in 1 minute
	http.SetCookie(w, &http.Cookie{
		Name:     "sso_cookie",
		Value:    s,
		Expires:  time.Now().Add(time.Minute * 2),
		HttpOnly: true,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
