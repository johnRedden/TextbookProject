package main

/*
filename.go by Allen J. Mills
    mm.d.yy

    Description
*/

import (
	"net/http"
	"time"
)

func ToCookie(res http.ResponseWriter, key string, value string, expiration time.Duration) {
	cookie := &http.Cookie{
		Name:   key,
		Value:  value,
		MaxAge: int(expiration.Seconds()),
	}
	http.SetCookie(res, cookie)
}

func FromCookie(req *http.Request, key string) (string, error) {
	cookie, err := req.Cookie(key)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
func UpdateCookie(res http.ResponseWriter, req *http.Request, key string, expiration time.Duration) error {
	cookie, err := req.Cookie(key)
	if err != nil {
		return err
	}
	cookie.MaxAge = int(expiration.Seconds())
	http.SetCookie(res, cookie)
	return nil
}
