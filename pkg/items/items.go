// Package items handles all the todo items
package items

import (
	"errors"
	"time"
)

const (
	// StatusDeleted is the status of the item returned after it's deleted
	StatusDeleted = "deleted"
	itemsBucket   = "items"
)

var (
	// ErrCreate is returned if there's an error in creating a new item
	ErrCreate = errors.New("Sorry, error creating an item")
	// ErrRead is returned if there's an error reading an item
	ErrRead = errors.New("Sorry, unable to fetch the item")
)

// Item holds a single item
type Item struct {
	// ID is the unique identifier of a single item
	ID string `json:"id,omitempty" bson:"_id,omitempty"`
	// Title is the title of a single item
	Title string `json:"title,omitempty" bson:"title,omitempty"`
	// Description is the description of a single item
	Description string `json:"description,omitempty" bson:"description,omitempty"`
	// Status is the current status of the item, it's set only while returning a deleted item
	Status string `json:"status,omitempty" bson:"status,omitempty"`
	// OwnerID is the unique identifier of an owner
	OwnerID string `json:"ownerID,omitempty" bson:"ownerID,omitempty"`
	// CreatedAt is a UTC timestamp of when the item was created
	CreatedAt *time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	// ModifiedAt is the UTC timestamp of when the item was last updated
	ModifiedAt *time.Time `json:"modifiedAt,omitempty" bson:"modifiedAt,omitempty"`
}

// Create creates a new item
func (s *Service) Create(data map[string]string) (*Item, error) {
	now := time.Now()
	item := Item{
		Title:       data["title"],
		Description: data["description"],
		OwnerID:     data["ownerID"],
		CreatedAt:   &now,
		ModifiedAt:  nil,
	}

	meta, err := s.store.Save(itemsBucket, item)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, ErrCreate
	}
	item.ID = meta.ID

	return &item, nil
}

// Read reads an item given the item ID
func (s *Service) Read(id string) (*Item, error) {
	item := Item{}
	_, err := s.store.FindOne(itemsBucket, id, &item)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}
	item.ID = id
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
	now := time.Now()
	item.ModifiedAt = &now
	item.ID = ""

	meta, err := s.store.UpdateOne(itemsBucket, data.ID, item)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}
	item.ID = meta.ID

	return item, nil
}

// Delete deletes an item given the ID
func (s *Service) Delete(id string) (*Item, error) {
	item, err := s.Read(id)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	meta, err := s.store.DeleteOne(itemsBucket, id)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	item.ID = meta.ID
	item.Status = StatusDeleted
	return item, nil
}

// List returns the list of items given the owner ID
func (s *Service) List(owner string) ([]Item, error) {
	return nil, nil
}
