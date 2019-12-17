package auth

import (
	"bufio"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"github.com/pkg/errors"
	"os"
	"strings"
)

const AuthUserKey = "user"

type AuthAdapter interface {
	Auth(ctx context.Context, c *Credentials) (*Result, error)
}

type Service struct {
	key *rsa.PrivateKey
	adapters []AuthAdapter
}

func New(db *sql.DB) *Service {

	adapters := []AuthAdapter{
		&Storage{db: db},
	}

	privateKey, err := getPrivateKey("data/key.pem")
	if err != nil {
		panic(err)
	}

	return &Service{
		key: privateKey,
		adapters: adapters,
	}
}

func getPrivateKey(filePath string) (*rsa.PrivateKey, error){
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return generatePrivateKey(filePath)
	}

	return importPrivateKey(filePath)
}

func importPrivateKey(filePath string) (*rsa.PrivateKey, error) {
	privateKeyFile, err := os.Open(filePath)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	pemfileinfo, _ := privateKeyFile.Stat()
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)
	buffer := bufio.NewReader(privateKeyFile)
	_, err = buffer.Read(pembytes)
	data, _ := pem.Decode([]byte(pembytes))
	privateKeyFile.Close()

	return x509.ParsePKCS1PrivateKey(data.Bytes)
}

func generatePrivateKey(filePath string) (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	pemPrivateFile, err := os.Create(filePath)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var pemPrivateBlock = &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	err = pem.Encode(pemPrivateFile, pemPrivateBlock)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = pemPrivateFile.Close()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return privateKey, nil
}

func (s *Service) IsAuthenticated() bool {
	return false
}

func (s *Service) Identity(ctx context.Context) (*User, error) {
	u := ctx.Value("user")
	if u != nil {
		usr, ok := u.(*User)
		if ok {
			return usr, nil
		}
	}

	a := ctx.Value("auth_token")
	if a != nil {
		token, ok := u.(string)
		if ok {
			return s.VerifySignature(ctx, token)
		}
	}

	return nil, nil
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

func (s *Service) Auth(ctx context.Context, c *Credentials) (string, error) {

	results := make([]*Result, 0, len(s.adapters))
	for _, a := range s.adapters {
		r, err := a.Auth(ctx, c)
		if err != nil {
			return "", err
		}

		if r != nil {
			results = append(results, r)
		}

		if r.IsValid() {
			return s.sign(r.Identity, &s.key.PublicKey)
		}
	}

	messages := make([]string, 0, len(results))
	for _, r := range results {
		messages = append(messages, r.Message)
	}
	return "", errors.New(strings.Join(messages, ", "))
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
