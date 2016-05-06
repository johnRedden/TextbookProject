package main

/*
AUTH_session_cookie.go by Allen J. Mills
    mm.d.yy

    Description
*/

import (
	"net/http"
	"time"
)

// Internal Function
// Description:
// This function will add a key:value pair into http:response with a life of time.Duration.
func ToCookie(res http.ResponseWriter, key string, value string, expiration time.Duration) {
	cookie := &http.Cookie{
		Name:   key,
		Value:  value,
		MaxAge: int(expiration.Seconds()),
	}
	http.SetCookie(res, cookie)
}

// Internal Function
// Description:
// This function will retrieve a value that may exist in key from request:cookie.
//
// Returns:
//      value(string) - Value of key
//      failure?(error) - If any errors occur they exist here.
func FromCookie(req *http.Request, key string) (string, error) {
	cookie, err := req.Cookie(key)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// Internal Function
// Description:
// This function will update a value with new time.Duration
//
// Returns:
//      failure?(error) - If any errors occur they exist here.
func UpdateCookie(res http.ResponseWriter, req *http.Request, key string, expiration time.Duration) error {
	cookie, err := req.Cookie(key)
	if err != nil {
		return err
	}
	cookie.MaxAge = int(expiration.Seconds())
	http.SetCookie(res, cookie)
	return nil
}

// Internal Function
// Description:
// This function will delete key from response:Cookie
func DeleteCookie(res http.ResponseWriter, key string) {
	cookie := &http.Cookie{
		Name:    key,
		Value:   "",
		MaxAge:  int(0),
		Expires: time.Now(),
	}
	http.SetCookie(res, cookie)
}
