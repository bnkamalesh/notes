// Package users handles user authentication, user items etc.
package users

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
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
	hasher = sha512.New()

	// ErrEmail is returned when the email address provided is wrong
	ErrEmail = errors.New("Invalid or no email provided")
	// ErrCreate is returned if there's an error creating new user
	ErrCreate = errors.New("Sorry, an error occurred while creating new user")
	// ErrInvPwd is returned if the password is invalid
	ErrInvPwd = errors.New("Sorry, invalid or no password provided")
	// ErrUsrNotExists is returned when trying to login with an non-registered email
	ErrUsrNotExists = errors.New("Sorry, there's no user registered with that email")
	// ErrUsrExists is returned when trying to create a user with the same email
	ErrUsrExists = errors.New("Sorry, user with that email already exists")
	// ErrNotAuthenticated is returned when the user is not authenticated and trying to perform
	// an action which requires authentication
	ErrNotAuthenticated = errors.New("Sorry, the user is not authenticated")
	// ErrUnauthorized is returned whenever the user tries to perform an unauthorized action
	ErrUnauthorized = errors.New("Sorry, you're not authorized to perform this action")
	// ErrMalformedCipher is returned when the cipher text is invalid and cannot be used
	ErrMalformedCipher = errors.New("malformed ciphertext")
)

func newUserID() string {
	return fmt.Sprintf("user|%s", uuid.New().String())
}

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
	hashSalt := uuid.New().String()
	now := time.Now()
	user := &User{
		ID:         newUserID(),
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
	Password          []byte     `json:"-" bson:"password,omitempty"`
	Salt              string     `json:"-" bson:"salt,omitempty"`
	AuthToken         string     `bson:"-" json:"authToken,omitempty"`
	EncryptedPassword []byte     `bson:"-" json:"-"`
	CreatedAt         *time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	ModifiedAt        *time.Time `json:"modifiedAt,omitempty" bson:"modifiedAt,omitempty"`
}

func (u *User) ownerID() (string, error) {
	password, err := u.passwordStr()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash(u.Email, password)), nil
}

// passwordStr returns the password of user in clear-text based on the authentication token
// This function will work only if the user is authenticated and has a valid authToken
// This is not saved anywhere, and is only used as the key for encryption and decryption of userdata
func (u *User) passwordStr() (string, error) {
	if u.AuthToken == "" {
		return "", ErrNotAuthenticated
	}

	key, err := u.encryptionKey(u.AuthToken, u.Salt)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	if len(u.EncryptedPassword) < gcm.NonceSize() {
		return "", ErrMalformedCipher
	}

	str, err := gcm.Open(nil,
		u.EncryptedPassword[:gcm.NonceSize()],
		u.EncryptedPassword[gcm.NonceSize():],
		nil,
	)
	if err != nil {
		return "", nil
	}
	return string(str), nil
}

// setEncryptedPassword sets the encrypted password in user struct field
func (u *User) setEncryptedPassword(password string) error {
	key, err := u.encryptionKey(u.AuthToken, u.Salt)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return err
	}

	u.EncryptedPassword = gcm.Seal(nonce, nonce, []byte(password), nil)
	return nil
}

func (u *User) encryptionKey(key, salt string) ([32]byte, error) {
	b := hash(key, salt)
	var bk [32]byte
	copy(bk[:], b[:32])
	return bk, nil
}

// hash accepts a string and returns a SHA512 hashed string
func hash(str string, salt string) []byte {
	return hasher.Sum([]byte(str + salt))
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

// Update reads a user given the email
func (s *Service) Update(user *User, data map[string]string) (*User, error) {
	name := strings.TrimSpace(data["name"])
	password := strings.TrimSpace(data["password"])

	if len(name) != 0 {
		user.Name = name
	}

	if len(password) != 0 {
		// pwdHash := hash(password, user.Salt)
		// savedPwdHash := user.Password
		// if pwdHash != savedPwdHash {
		// Decrypt and encrypt all items of the user with the new password
		// Should update owner ID also
		// }
	}
	return user, nil
}

// Delete deletes the provided User
func (s *Service) Delete(user *User) (*User, error) {
	err := s.store.Delete(userBucket, map[string]interface{}{
		"id": user.ID,
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

// CreateItem adds a new item owned by the user
func (s *Service) CreateItem(user *User, data map[string]string) (*items.Item, error) {
	ownerID, err := user.ownerID()
	if err != nil {
		return nil, err
	}

	item, err := items.New(data, ownerID)
	if err != nil {
		return nil, err
	}
	pwd, err := user.passwordStr()
	if err != nil {
		return nil, err
	}

	key, err := user.encryptionKey(pwd, user.AuthToken)
	if err != nil {
		return nil, err
	}

	err = item.Encrypt(key)
	if err != nil {
		return nil, err
	}

	return s.items.Create(*item)
}

// UpdateItem updates an item owned by the user
func (s *Service) UpdateItem(user *User, itemID string, data map[string]string) (*items.Item, error) {
	ownerID, err := user.ownerID()
	if err != nil {
		return nil, err
	}

	item, err := s.items.Read(itemID)
	if err != nil {
		return nil, err
	}

	if item.OwnerID != ownerID {
		return nil, ErrUnauthorized
	}

	updatedItem, err := items.New(data, ownerID)
	if err != nil {
		return nil, err
	}
	updatedItem.ID = itemID

	pwd, err := user.passwordStr()
	if err != nil {
		return nil, err
	}

	key, err := user.encryptionKey(pwd, user.AuthToken)
	if err != nil {
		return nil, err
	}

	err = updatedItem.Encrypt(key)
	if err != nil {
		return nil, err
	}

	updatedItem, err = s.items.Update(itemID, *updatedItem)
	if err != nil {
		return nil, err
	}

	return updatedItem, err
}

// DeleteItem removes an item owned by the user
func (s *Service) DeleteItem(user *User, itemID string) (*items.Item, error) {
	ownerID, err := user.ownerID()
	if err != nil {
		return nil, err
	}

	item, err := s.items.Read(itemID)
	if err != nil {
		return nil, err
	}

	if item.OwnerID != ownerID {
		return nil, ErrUnauthorized
	}

	item, err = s.items.Delete(itemID)
	return item, err
}

// Items returns list of items the user owns
func (s *Service) Items(user *User, start, limit int) ([]items.Item, error) {
	ownerID, err := user.ownerID()
	if err != nil {
		return nil, err
	}
	ii, err := s.items.List(ownerID, start, limit)
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

	pwd, err := user.passwordStr()
	if err != nil {
		return nil, err
	}

	key, err := user.encryptionKey(pwd, user.AuthToken)
	if err != nil {
		return nil, err
	}

	err = i.Decrypt(key)
	if err != nil {
		return nil, err
	}
	return i, nil
}
