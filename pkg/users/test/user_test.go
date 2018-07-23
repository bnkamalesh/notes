package users

import (
	"testing"

	"github.com/bnkamalesh/notes/pkg/items"

	"github.com/bnkamalesh/notes/pkg/platform/logger"
	"github.com/bnkamalesh/notes/pkg/platform/storage"
	"github.com/bnkamalesh/notes/pkg/users"
)

func service() (*users.Service, error) {
	store, err := storage.New(storage.Config{
		Name:  "gonotes_test",
		Hosts: []string{"127.0.0.1:27017"},
	})
	if err != nil {
		return nil, err
	}
	logHandler := logger.New()
	iS := items.NewService(store, logHandler)
	service := users.NewService(store, logHandler, iS)
	return &service, nil
}

func newUser() (*users.User, map[string]string, error) {
	payload := map[string]string{
		"name":     "John Smith",
		"email":    "jsmith@example.com",
		"password": "hello world",
	}
	item, err := users.New(payload)
	return item, payload, err
}
func TestCreate(t *testing.T) {
	s, err := service()
	if err != nil {
		t.Fatal(err.Error())
	}
	u, _, err := newUser()
	if err != nil {
		t.Fatal(err.Error())
	}

	createdUsr, err := s.Create(*u)
	if err != nil {
		t.Fatal(err.Error())
	}
	if createdUsr.ID == "" {
		t.Fatal("Invalid or no ID")
	}
	if createdUsr.Email != u.Email {
		t.Fatalf("Expected email '%s', got '%s'", u.Email, createdUsr.Email)
	}
	if createdUsr.Name != u.Name {
		t.Fatalf("Expected name '%s', got '%s'", u.Name, createdUsr.Name)
	}
}

func TestAuth(t *testing.T) {
	s, err := service()
	if err != nil {
		t.Fatal(err.Error())
	}
	u, payload, err := newUser()
	if err != nil {
		t.Fatal(err.Error())
	}

	createdUsr, err := s.Create(*u)
	if err != nil {
		t.Fatal(err.Error())
		return
	}

	authUser, err := s.Authenticate(createdUsr.Email, payload["password"])
	if err != nil {
		t.Fatalf("authenticate failed, email '%s',  password '%s', error: '%s'", createdUsr.Email, payload["password"], err.Error())
	}
	if createdUsr.ID != authUser.ID {
		t.Fatalf("Expected user ID, '%s', got '%s'", createdUsr.ID, authUser.ID)
	}
}
func TestAddItem(t *testing.T) {
	s, err := service()
	if err != nil {
		t.Fatal(err.Error())
	}
	u, payload, err := newUser()
	if err != nil {
		t.Fatal(err.Error())
	}

	createdUsr, err := s.Create(*u)
	if err != nil {
		t.Fatal(err.Error())
		return
	}

	authUser, err := s.Authenticate(createdUsr.Email, payload["password"])
	if err != nil {
		t.Fatalf("%s %s %s", createdUsr.Email, payload["password"], err.Error())
	}
	if createdUsr.ID != authUser.ID {
		t.Fatalf("Expected ID, '%s', got '%s'", createdUsr.ID, authUser.ID)
	}

	itemPayload := map[string]string{
		"title":       "Hello",
		"description": "well well well",
	}

	item, err := s.AddItem(authUser, itemPayload)
	if err != nil {
		t.Fatal(err.Error())
	}

	rI, err := s.Item(authUser, item.ID)
	if err != nil {
		t.Fatal(err.Error())
	}
	if rI.ID != item.ID {
		t.Fatalf("Expected item ID '%s', got '%s'", item.ID, rI.ID)
	}

	if rI.Description != itemPayload["description"] {
		t.Fatalf("Expected item description '%s', got '%s'", itemPayload["description"], rI.Description)
	}
}
