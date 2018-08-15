package items

import (
	"testing"

	"github.com/bnkamalesh/notes/pkg/platform/logger"
	"github.com/bnkamalesh/notes/pkg/platform/storage"
)

func service() (*Service, error) {
	store, err := storage.New(storage.Config{
		Name:  "gonotes_test",
		Hosts: []string{"127.0.0.1:27017"},
	})
	if err != nil {
		return nil, err
	}
	logHandler := logger.New()
	service := NewService(store, logHandler)
	return &service, nil
}

func newItem() (*Item, map[string]string, error) {
	payload := map[string]string{
		"title":       "Hello world",
		"description": "Hello world description",
	}
	item, err := New(payload, "testOwner")
	return item, payload, err
}

func TestCreate(t *testing.T) {
	s, err := service()
	if err != nil {
		t.Fatal(err.Error())
	}

	item, payload, err := newItem()
	if err != nil {
		t.Fatal(err.Error())
	}
	// Test create
	item, err = s.Create(*item)
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
	if item.OwnerID != "testOwner" {
		t.Fatalf("Invalid OwnerID, got '%s' expected '%s'", item.OwnerID, "testOwner")
	}
	_, err = s.Delete(item.ID)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestRead(t *testing.T) {
	s, err := service()
	if err != nil {
		t.Fatal(err.Error())
	}
	item, _, err := newItem()
	if err != nil {
		t.Fatal(err.Error())
	}
	item, err = s.Create(*item)
	if err != nil {
		t.Fatal(err.Error())
	}

	itemFromDB, err := s.Read(item.ID)
	if err != nil {
		t.Fatal(err.Error())
	}

	if itemFromDB.ID != item.ID {
		t.Fatalf("Invalid ID, got '%s' expected '%s'", itemFromDB.ID, item.ID)
	}
	_, err = s.Delete(itemFromDB.ID)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestUpdate(t *testing.T) {
	s, err := service()
	if err != nil {
		t.Fatal(err.Error())
	}
	item, _, err := newItem()
	if err != nil {
		t.Fatal(err.Error())
	}
	item, err = s.Create(*item)
	if err != nil {
		t.Fatal(err.Error())
	}

	itemFromDB, err := s.Read(item.ID)
	if err != nil {
		t.Fatal(err.Error())
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
	_, err = s.Delete(itemFromDB.ID)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestDelete(t *testing.T) {
	s, err := service()
	if err != nil {
		t.Fatal(err.Error())
	}
	item, _, err := newItem()
	if err != nil {
		t.Fatal(err.Error())
	}
	item, err = s.Create(*item)
	if err != nil {
		t.Fatal(err.Error())
	}

	itemFromDB, err := s.Read(item.ID)
	if err != nil {
		t.Fatal(err.Error())
	}

	deletedItem, err := s.Delete(itemFromDB.ID)
	if err != nil {
		t.Fatal(err.Error())
	}

	if deletedItem.ID != itemFromDB.ID {
		t.Fatalf("Invalid ID, got '%s' expected '%s'", deletedItem.ID, itemFromDB.ID)
	}
}

func TestList(t *testing.T) {
	s, err := service()
	if err != nil {
		t.Fatal(err.Error())
	}
	payload1 := map[string]string{
		"title":       "Hello world",
		"description": "Hello world description",
	}
	item1, err := New(payload1, "testOwner")
	if err != nil {
		t.Fatal(err.Error())
	}

	cItem1, err := s.Create(*item1)
	if err != nil {
		t.Fatal(err.Error())
	}

	payload2 := map[string]string{
		"title":       "Hello world 2",
		"description": "Hello world description 2",
	}
	item2, err := New(payload2, "testOwner")
	if err != nil {
		t.Fatal(err.Error())
	}

	cItem2, err := s.Create(*item2)
	if err != nil {
		t.Fatal(err.Error())
	}

	ii, err := s.List("testOwner", 0, 100)
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(ii) != 2 {
		t.Fatalf("Expected '%d', got '%d' results.", 2, len(ii))
	} else {
		for _, i := range ii {
			if i.Title != cItem1.Title && i.Title != cItem2.Title {
				t.Fatalf("Expected '%s' or '%s', got '%s'", cItem1.Title, cItem2.Title, i.Title)
			}
			id := i.ID
			_, err := s.Delete(id)
			if err != nil {
				t.Fatal(err)
			}
		}
	}

}
