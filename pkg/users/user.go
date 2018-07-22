// Package users handles user authentication, user items etc.
package users

import (
	"crypto/aes"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/bnkamalesh/notes/pkg/items"
	"github.com/bnkamalesh/notes/pkg/platform/storage"
)

const (
	userBucket = "users"
)

var (
	hashSalt = uuid.New().String()
	hasher   = sha512.New()

	// ErrEmail is returned when the email address provided is wrong
	ErrEmail = errors.New("Invalid or no email provided")
	// ErrCreate is returned if there's an error creating new user
	ErrCreate = errors.New("Sorry, an error occurred while creating new user")
	// ErrInvPwd is returned if the password is invalid
	ErrInvPwd = errors.New("Sorry, invalid or no password provided")
	// ErrUsrNotExists is returned when trying to login with an non-registered email
	ErrUsrNotExists = errors.New("Sorry, there's no user with that email")
	// ErrUsrExists is returned when trying to create a user with the same email
	ErrUsrExists = errors.New("Sorry, user with that email already exists")
	// ErrNotAuthenticated is returned when the user is not authenticated and trying to perform
	// an action which requires authentication
	ErrNotAuthenticated = errors.New("Sorry, the user is not authenticated")
)

// New returns a user instance based on the provided data
func New(data map[string]string) (*User, error) {
	email := strings.TrimSpace(data["email"])
	password := data["password"]
	if email == "" {
		return nil, ErrEmail
	}
	if password == "" {
		return nil, ErrInvPwd
	}

	now := time.Now()
	user := &User{
		ID:         uuid.New().String(),
		Name:       data["name"],
		Email:      email,
		Salt:       hashSalt,
		Password:   hash(password, hashSalt),
		CreatedAt:  &now,
		ModifiedAt: nil,
	}
	return user, nil
}

// User struct holds all the user details
type User struct {
	ID                string     `json:"id,omitempty" bson:"id,omitempty"`
	Name              string     `json:"name,omitempty" bson:"name,omitempty"`
	Email             string     `json:"email,omitempty" bson:"email,omitempty"`
	Password          string     `json:"-" bson:"password,omitempty"`
	Salt              string     `json:"-" bson:"salt,omitempty"`
	authToken         string     `bson:"-"`
	encryptedPassword string     `bson:"-"`
	ownerID           string     `bson:"-"`
	CreatedAt         *time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	ModifiedAt        *time.Time `json:"modifiedAt,omitempty" bson:"modifiedAt,omitempty"`
}

// passwordStr returns the password of user in clear-text based on the authentication token
// This function will work only if the user is authenticated and has a valid authToken
// This is not saved anywhere, and is only used as the key for encryption and decryption of userdata
func (u *User) passwordStr() (string, error) {
	if u.authToken == "" {
		return "", ErrNotAuthenticated
	}
	block, err := aes.NewCipher([]byte(u.authToken))
	if err != nil {
		return "", err
	}
	dst := make([]byte, 0)
	block.Decrypt(dst, []byte(u.encryptedPassword))
	return string(dst), nil
}

// hash accepts a string and returns a SHA512 hashed string
func hash(str string, salt string) string {
	hasher.Write([]byte(str + salt))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha
}

// Create creates a new user
func (s *Service) Create(user User) (*User, error) {
	email := user.Email
	_, err := s.Read(email)
	if err != nil {
		if err != ErrUsrNotExists {
			s.logger.Fatal(err.Error())
			return nil, err
		}
	} else {
		return nil, ErrUsrExists
	}

	_, err = s.store.Save(userBucket, user)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, ErrCreate
	}

	return &user, nil
}

// Read reads a user given the email
func (s *Service) Read(email string) (*User, error) {
	user := User{}
	_, err := s.store.FindOne(
		userBucket,
		map[string]interface{}{
			"email": email,
		},
		nil,
		nil,
		&user)
	if err != nil {
		if err == storage.ErrNotFound {
			return nil, ErrUsrNotExists
		}
		s.logger.Error(err.Error())
		return nil, err
	}
	return &user, nil
}

// AddItem adds a new item owned by the user
func (s *Service) AddItem(user *User, data map[string]string) (*items.Item, error) {
	item, err := items.New(data, user.ownerID)
	if err != nil {
		return nil, err
	}
	key, err := user.passwordStr()
	if err != nil {
		return nil, err
	}

	err = item.Encrypt(key)
	if err != nil {
		return nil, err
	}

	return s.items.Create(*item)
}

// Items returns list of items the user owns
func (s *Service) Items(user *User, start, limit int) ([]items.Item, error) {
	ii, err := s.items.List(user.ownerID, start, limit)
	if err != nil {
		return nil, err
	}
	return ii, nil
}

// Item returns a decrypted item
func (s *Service) Item(user *User, itemID string) (*items.Item, error) {
	i, err := s.items.Read(itemID)
	if err != nil {
		return nil, err
	}

	key, err := user.passwordStr()
	if err != nil {
		return nil, err
	}

	i.Decrypt(key)

	return i, nil
}