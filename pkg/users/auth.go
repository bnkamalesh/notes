package users

import (
	"crypto/aes"
	"encoding/hex"

	"github.com/google/uuid"
)

// authToken generates an auth token for the user
func authToken(user *User) string {
	return hex.EncodeToString(hash(uuid.New().String(), user.Email))
}

// Authenticate authenticates a user and returns the user instance along with the auth token
func (s *Service) Authenticate(email, password string) (*User, error) {
	user, err := s.Read(email)
	if err != nil {
		return nil, err
	}

	pwdHash := hex.EncodeToString(hash(password, user.Salt))
	savedPwdHash := hex.EncodeToString(user.Password)
	if pwdHash != savedPwdHash {
		return nil, ErrInvPwd
	}

	user.authToken = authToken(user)
	block, err := aes.NewCipher([]byte(user.authToken))
	if err != nil {
		return nil, err
	}

	src := []byte(password)
	dst := make([]byte, len(src))
	block.Encrypt(dst, src)
	user.encryptedPassword = string(dst)
	user.ownerID = hex.EncodeToString(hash(email, password))

	return user, nil
}

// AuthUser returns an authenticated user instance from the auth token
func (s *Service) AuthUser(authToken string) (*User, error) {
	user := User{}
	return &user, nil
}
