package items

import (
	"testing"

	"github.com/bnkamalesh/gonotes/pkg/platform/logger"
	"github.com/bnkamalesh/gonotes/pkg/platform/storage"
)

func service() (*Service, error) {
	store, err := storage.New(storage.Config{
		Name:  "gonotes",
		Hosts: []string{"127.0.0.1:27017"},
	})
	if err != nil {
		return nil, err
	}
	logHandler := logger.New()
	service := NewService(store, logHandler)
	return &service, nil
}

func TestCRUD(t *testing.T) {
	s, err := service()
	if err != nil {
		t.Fatal(err.Error())
	}

	payload := map[string]string{
		"title":       "Hello world",
		"description": "Hello world description",
		"ownerID":     "testOwner",
	}

	// Test create
	item, err := s.Create(payload)
	if err != nil {
		t.Fatal(err.Error())
	}

	if item.ID == "" {
		t.Fatal("Invalid ID")
	}
	if item.Title != payload["title"] {
		t.Fatalf("Invalid title, got '%s' expected '%s'", item.Title, payload["title"])
	}
	if item.Description != payload["description"] {
		t.Fatalf("Invalid description, got '%s' expected '%s'", item.Description, payload["description"])
	}
	if item.OwnerID != payload["ownerID"] {
		t.Fatalf("Invalid OwnerID, got '%s' expected '%s'", item.OwnerID, payload["ownerID"])
	}

	// Test read
	itemFromDB, err := s.Read(item.ID)
	if err != nil {
		t.Fatal(err.Error())
	}

	if itemFromDB.ID != item.ID {
		t.Fatalf("Invalid ID, got '%s' expected '%s'", itemFromDB.ID, item.ID)
	}

	// Test update
	const updateTitle = "updated title, hello"
	itemFromDB.Title = updateTitle
	updatedItem, err := s.Update(itemFromDB.ID, *itemFromDB)
	if err != nil {
		t.Fatal(err.Error())
	}
	if updatedItem.ID != itemFromDB.ID {
		t.Fatalf("Invalid ID, got '%s' expected '%s'", updatedItem.ID, itemFromDB.ID)
	}

	if updatedItem.Title != updateTitle {
		t.Fatalf("Invalid title, got '%s' expected '%s'", updatedItem.Title, updateTitle)
	}

	// Test delete
	deletedItem, err := s.Delete(updatedItem.ID)
	if err != nil {
		t.Fatal(err.Error())
	}

	if deletedItem.ID != updatedItem.ID {
		t.Fatalf("Invalid ID, got '%s' expected '%s'", deletedItem.ID, updatedItem.ID)
	}
}
