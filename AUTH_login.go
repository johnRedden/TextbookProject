package main

/*
filename.go by Allen J. Mills
    mm.d.yy

    Description
*/
import (
	"errors"
	"github.com/Esseh/retrievable"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

var (
	ErrUserExists = errors.New("Login Error: This user already exists!")
)

type LoginAccount struct {
	UserKey  int64
	Password []byte
}

func (l *LoginAccount) Key(ctx context.Context, id interface{}) *datastore.Key {
	return datastore.NewKey(ctx, "Logins", id.(string), 0, nil)
}

func GetUserIDFromLogin(ctx context.Context, email, password string) (int64, error) {
	uLogin := LoginAccount{}
	if getErr := retrievable.GetEntity(ctx, &uLogin, email); getErr != nil {
		return -1, getErr
	}

	if compareErr := bcrypt.CompareHashAndPassword(uLogin.Password, []byte(password)); compareErr != nil {
		return -1, compareErr
	}

	return uLogin.UserKey, nil
}

func CreateUserFromLogin(ctx context.Context, email, password string, u *User) (*User, error) {
	// Step 1: Verify user does not exist
	if checkLoginErr := retrievable.GetEntity(ctx, &LoginAccount{}, email); checkLoginErr == nil {
		return u, ErrUserExists
	} else if checkLoginErr != datastore.ErrNoSuchEntity {
		return u, checkLoginErr
	}

	// Step 2: Place user into datastore
	uKey, putUserErr := retrievable.PlaceInDatastore(ctx, int64(0), u)
	if putUserErr != nil {
		return u, putUserErr
	}

	if u.ID == int64(0) {
		return u, errors.New("HEY, DATASTORE IS STUPID")
	}

	// Step 3: Encrypt user password
	cPass, cErr := bcrypt.GenerateFromPassword([]byte(password), 0)
	if cErr != nil {
		return u, cErr
	}

	// Step 4: Create user login profile
	uLogin := LoginAccount{}
	uLogin.Password = cPass
	uLogin.UserKey = uKey.IntID()
	_, putErr := retrievable.PlaceEntity(ctx, email, &uLogin)
	return u, putErr
}

func DeleteUserIDAndLoginKey(ctx context.Context, email, password string) error {
	uId, loginErr := GetUserIDFromLogin(ctx, email, password)
	if loginErr != nil {
		return loginErr
	}

	if uDeleteErr := retrievable.DeleteEntity(ctx, (&User{}).Key(ctx, uId)); uDeleteErr != nil {
		return uDeleteErr
	}

	return retrievable.DeleteEntity(ctx, (&LoginAccount{}).Key(ctx, email))
}

func ValidatePassword() (bool, error) {
	return true, ErrNotImplemented
}
