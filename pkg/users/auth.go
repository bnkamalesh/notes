package users

import (
	"bytes"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
)

// setAuthCache will store the user object in cache
func (s *Service) setAuthCache(token string, user *User) error {
	return s.cache.Set(token, user, time.Hour*24)
}

func (s *Service) getAuthCache(token string) (*User, error) {
	user := User{}
	err := s.cache.Get(token, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// cacheAuthToken will store the authtoken in cache after mixing it with salt
func cacheAuthToken(token string, salt string) string {
	return string(hash(token, salt))
}

// authToken generates an auth token for the user
func authToken(user *User) string {
	return hex.EncodeToString(hash(uuid.New().String(), user.Email))
}

// Authenticate authenticates a user and returns the user instance along with the auth token
func (s *Service) Authenticate(email, password, tokenSalt string) (*User, error) {
	user, err := s.Read(email)
	if err != nil {
		return nil, err
	}

	pwdHash := hash(password, user.Salt)
	savedPwdHash := user.Password
	if !bytes.Equal(pwdHash, savedPwdHash) {
		return nil, ErrInvPwd
	}

	user.AuthToken = authToken(user)
	err = user.setEncryptedPassword(password)
	if err != nil {
		return nil, err
	}
	err = s.setAuthCache(cacheAuthToken(user.AuthToken, tokenSalt), user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// AuthUser returns an authenticated user instance from the auth token
func (s *Service) AuthUser(authToken string, tokenSalt string) (*User, error) {
	return s.getAuthCache(cacheAuthToken(authToken, tokenSalt))
}
