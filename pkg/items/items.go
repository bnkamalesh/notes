// Package items handles all the todo items
package items

import "time"

const (
	itemsBucket = "items"
)

// Item holds a single item
type Item struct {
	// ID is the unique identifier of a single item
	ID string `json:"id,omitempty" bson:"_id,omitempty"`
	// Title is the title of a single item
	Title string `json:"title,omitempty" bson:"title,omitempty"`
	// Description is the description of a single item
	Description string `json:"description,omitempty" bson:"description,omitempty"`
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
		return nil, err
	}
	item.ID = meta.ID

	return &item, nil
}

// Read reads an item given the item ID
func (s *Service) Read(id string) (*Item, error) {
	item := Item{}
	_, err := s.store.FindOne(itemsBucket, id, &item)
	if err != nil {
		return nil, err
	}
	item.ID = id
	return &item, nil
}

// Update updates an item given the ID
func (s *Service) Update(id string, data Item) (*Item, error) {
	item, err := s.Read(id)
	if err != nil {
		return nil, err
	}

	item.Title = data.Title
	item.Description = data.Description
	now := time.Now()
	item.ModifiedAt = &now
	item.ID = ""

	meta, err := s.store.UpdateOne(itemsBucket, data.ID, item)
	if err != nil {
		return nil, err
	}
	item.ID = meta.ID

	return item, nil
}

// Delete deletes an item given the ID
func (s *Service) Delete(id string) (*Item, error) {
	item, err := s.Read(id)
	if err != nil {
		return nil, err
	}
	meta, err := s.store.DeleteOne(itemsBucket, id)
	if err != nil {
		return nil, err
	}
	item.ID = meta.ID
	return item, nil
}

// List returns the list of items given the owner ID
func (s *Service) List(owner string) ([]Item, error) {
	return nil, nil
}
