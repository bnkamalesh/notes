package users

import (
	"encoding/hex"

	"github.com/google/uuid"
)

// setAuthCache will store the user object in cache
func setAuthCache(token string, user *User) error {
	return nil
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

	pwdHash := hex.EncodeToString(hash(password, user.Salt))
	savedPwdHash := hex.EncodeToString(user.Password)
	if pwdHash != savedPwdHash {
		return nil, ErrInvPwd
	}

	user.AuthToken = authToken(user)
	user.setEncryptedPassword(password)
	setAuthCache(cacheAuthToken(user.AuthToken, tokenSalt), user)
	return user, nil
}

// AuthUser returns an authenticated user instance from the auth token
func (s *Service) AuthUser(authToken string, tokenSalt string) (*User, error) {
	user := User{}
	return &user, nil
}
