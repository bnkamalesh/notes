// Package items handles all the todo items
package items

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	// StatusDeleted is the status of the item returned after it's deleted
	StatusDeleted = "deleted"
	itemsBucket   = "items"
	minStart      = 0
	maxLimit      = 50
)

var (
	// ErrCreate is returned if there's an error in creating a new item
	ErrCreate = errors.New("Sorry, an error occurred while creating")
	// ErrRead is returned if there's an error reading an item
	ErrRead = errors.New("Sorry, unable to fetch item")
	// ErrInvOwnerID is returned if the owner ID is blank or invalid
	ErrInvOwnerID = errors.New("Sorry, invalid owner ID provided")
)

// Item holds a single item
type Item struct {
	// ID is the unique identifier of a single item
	ID string `json:"id,omitempty" bson:"id,omitempty"`
	// Title is the title of a single item
	Title string `json:"title,omitempty" bson:"title,omitempty"`
	// Description is the description of a single item
	Description string `json:"description,omitempty" bson:"description,omitempty"`
	// Status is the current status of the item, it's set only while returning a deleted item
	Status string `json:"status,omitempty" bson:"status,omitempty"`
	// OwnerID is the unique identifier of an owner
	OwnerID string `json:"-" bson:"ownerID,omitempty"`
	// Blob stores the encrypted byte of Item
	Blob []byte `json:"-" bson:"blob,omitempty"`
	// CreatedAt is a UTC timestamp of when the item was created
	CreatedAt *time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	// ModifiedAt is the UTC timestamp of when the item was last updated
	ModifiedAt *time.Time `json:"modifiedAt,omitempty" bson:"modifiedAt,omitempty"`
}

func newItemID() string {
	return fmt.Sprintf("item_%s", uuid.New().String())
}

// New returns a new instance of Item with the provided data
func New(data map[string]string, ownerID string) (*Item, error) {
	ownerID = strings.TrimSpace(ownerID)
	if ownerID == "" {
		return nil, ErrInvOwnerID
	}

	now := time.Now()
	return &Item{
		ID:          newItemID(),
		Title:       strings.TrimSpace(data["title"]),
		Description: strings.TrimSpace(data["description"]),
		OwnerID:     ownerID,
		CreatedAt:   &now,
		ModifiedAt:  &now,
	}, nil
}

// Encrypt encrypts the item description and sets the Blob with encrypted bytes
func (i *Item) Encrypt(key [32]byte) error {
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

	i.Blob = gcm.Seal(nonce, nonce, []byte(i.Description), nil)

	// Emptying the description to prevent it from being saved as plain text
	i.Description = ""
	return nil
}

// Decrypt decrpyts an item Description with the provided key
func (i *Item) Decrypt(key [32]byte) error {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	if len(i.Blob) < gcm.NonceSize() {
		return errors.New("malformed ciphertext")
	}

	str, err := gcm.Open(nil,
		i.Blob[:gcm.NonceSize()],
		i.Blob[gcm.NonceSize():],
		nil,
	)
	if err != nil {
		return err
	}
	i.Description = string(str)
	return nil
}

// Create creates a new item
func (s *Service) Create(item Item) (*Item, error) {
	_, err := s.store.Save(itemsBucket, item)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, ErrCreate
	}
	return &item, nil
}

// Read reads an item given the item ID
func (s *Service) Read(id string) (*Item, error) {
	item := Item{}
	_, err := s.store.FindOne(
		itemsBucket,
		map[string]interface{}{"id": id},
		nil,
		[]string{"-modifiedAt"},
		&item)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	return &item, nil
}

// Update updates an item given the ID
func (s *Service) Update(id string, data Item) (*Item, error) {
	item, err := s.Read(id)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	item.Title = data.Title
	item.Description = data.Description
	item.Blob = data.Blob

	now := time.Now()
	item.ModifiedAt = &now

	err = s.store.Update(
		itemsBucket,
		map[string]interface{}{
			"id": data.ID,
		}, item)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	return item, nil
}

// Delete deletes an item given the ID
func (s *Service) Delete(id string) (*Item, error) {
	item, err := s.Read(id)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	err = s.store.Delete(
		itemsBucket,
		map[string]interface{}{
			"id": id,
		})
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	item.Status = StatusDeleted
	return item, nil
}

// List returns the list of items given the owner ID
func (s *Service) List(ownerID string, start, limit int) ([]Item, error) {
	query := map[string]interface{}{
		"ownerID": ownerID,
	}

	if start < minStart {
		start = minStart
	}

	if limit <= 0 || limit > maxLimit {
		limit = maxLimit
	}

	out := make([]Item, 0)
	_, err := s.store.Find(itemsBucket, query, nil, []string{"-modifiedAt"}, start, limit, &out)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	return out, nil
}
