package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtCustomClaims struct {
	jwt.RegisteredClaims
}

func main() {
	key := []byte("secret")

	s := http.NewServeMux()
	s.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		claims := &jwtCustomClaims{
			jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 1)),
				Issuer:    "haproxy",
			},
		}
		t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		s, err := t.SignedString(key)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("error: %s\n", err)))
			return
		}
		// set cookie that expires in 1 minute
		http.SetCookie(w, &http.Cookie{
			Name:     "sso_cookie",
			Value:    s,
			Expires:  time.Now().Add(time.Minute * 2),
			HttpOnly: true,
		})
		w.Write([]byte(s))
		w.Write([]byte("\n"))
		w.Write([]byte(fmt.Sprintf("alg: %s\n", t.Method.Alg())))
		w.Write([]byte(fmt.Sprintf("iss: %s\n", claims.Issuer)))
		w.Write([]byte(fmt.Sprintf("exp: %s\n", claims.ExpiresAt.Time.String())))
	})
	http.ListenAndServe(":8083", s)
}
