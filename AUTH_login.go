package main

/*
filename.go by Allen J. Mills
    mm.d.yy

    Description
*/
import (
	"errors"
	"github.com/Esseh/retrievable"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

var (
	ErrUserExists = errors.New("Register Error: This user already exists!")
)

type UserLogin struct {
	UID int64
}

func (ul *UserLogin) Key(ctx context.Context, id interface{}) *datastore.Key {
	return datastore.NewKey(ctx, "Logins", id.(string), 0, nil)
}

func CreateLogin(ctx context.Context, u *User) error {
	_, lgnErr := retrievable.PlaceInDatastore(ctx, u.Email, &UserLogin{u.ID})
	return lgnErr
}

func GetUIDFromLogin(ctx context.Context, email string) (int64, error) {
	ul := &UserLogin{}
	getErr := retrievable.GetFromDatastore(ctx, email, ul)
	return ul.UID, getErr
}
