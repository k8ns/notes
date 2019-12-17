package auth

import (
	"context"
	"errors"
)

type User struct {
	Id uint
	Identity string
	Password string
}

func GetUser(ctx context.Context) (*User, error) {
	u, ok := ctx.Value(AuthUserKey).(*User)
	if !ok {
		return nil, errors.New("unauthorized")
	}
	return u, nil
}
