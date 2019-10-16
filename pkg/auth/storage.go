package auth

import (
	"database/sql"
	"errors"
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

