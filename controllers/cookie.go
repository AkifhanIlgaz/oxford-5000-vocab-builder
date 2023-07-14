package controllers

import "net/http"

const (
	CookieSession = "session"
)

func newCookie(name, value string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		HttpOnly: true,
		Path:     "/",
	}
}

func setCookie(w http.ResponseWriter, name, value string) {
	http.SetCookie(w, newCookie(name, value))
}
