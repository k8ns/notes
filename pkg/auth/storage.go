package auth

import (
	"context"
	"database/sql"
	"errors"

    "golang.org/x/crypto/bcrypt"
)

type Storage struct {
	db *sql.DB
}

func (s *Storage) Find(identity string) (*User, error) {

	rows, err := s.db.Query("SELECT * FROM users WHERE identity = ?", identity)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		u := &User{}
		err = rows.Scan(&u.Id, &u.Identity, &u.Password)
		return u, err
	}
	return nil, errors.New("not found")
}


func (s *Storage) Auth(ctx context.Context, c *Credentials) (*Result, error) {
	u, err := s.Find(c.Identity)
	if err != nil {
		return nil, err
	}

	if u == nil {
		return &Result{Code: ResultNoIdentity, Message: "no user with such email"}, nil
	}

	if u.Password == "" {
		return &Result{Code: ResultWrongCredentials, Message: "no password"}, nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(c.Credentials))
	if err != nil {
		//encrypted, _ := bcrypt.GenerateFromPassword([]byte(c.Credentials), bcrypt.DefaultCost)
		return &Result{
			Code: ResultWrongCredentials,
			Message: "wrong credentials",
		}, nil
	}

	return &Result{Code: ResultSuccess, Message: "Access grunted", Identity: u}, nil
}

