package users

import (
	"crypto/aes"
	"fmt"

	"github.com/google/uuid"
)

// authToken generates an auth token for the user
func authToken(user *User) string {
	return hash(uuid.New().String(), user.Password)
}

// Authenticate authenticates a user and returns the user instance along with the auth token
func (s *Service) Authenticate(email, password string) (*User, error) {
	user, err := s.Read(email)
	if err != nil {
		return nil, err
	}

	pwdHash := hash(password, user.Salt)
	if pwdHash != user.Password {
		fmt.Println("pwdHash:", pwdHash, " user.Pwd:", user.Password, " salt:", user.Salt)
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
	user.ownerID = hash(email, password)
	return user, nil
}

// AuthUser returns an authenticated user instance from the auth token
func (s *Service) AuthUser(authToken string) (*User, error) {
	user := User{}
	return &user, nil
}
