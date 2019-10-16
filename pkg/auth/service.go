package auth

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
)

const AuthUserKey = "user"

type Service struct {
	storage *Storage
	key *rsa.PrivateKey
}

func New(cfg *Config, db *sql.DB) *Service {

	k, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	return &Service{
		storage: &Storage{db: db},
		key: k,
	}
}

func (s *Service) Auth(ctx context.Context, c *Credential) (string, error) {
	u, err := s.storage.Find(c.Identity)
	if err != nil {
		return "", errors.New("denied")
	}

	if !s.VerifyPassword(c.Credential, u.Password) {
		return "", errors.New("denied")
	}

	return s.sign(u, &s.key.PublicKey)
}

func (s *Service) sign(u *User, key *rsa.PublicKey) (string, error) {

	secretMessage, err := json.Marshal(u)
	if err != nil {
		return "", err
	}

	label := []byte("A")

	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, key, secretMessage, label)
	if err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(ciphertext)

	return encoded, nil
}

func (s *Service) VerifySignature(ctx context.Context, ciphertext string) (*User, error) {
	decoded, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}

	label := []byte("A")
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, s.key, decoded, label)
	if err != nil {
		return nil, err
	}

    u := &User{}
	err = json.Unmarshal(plaintext, u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Service) VerifyPassword(c, p string) bool {
	// todo: encrypt password
	return c == p
}
