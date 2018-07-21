// Package storage handles the primary data store
package storage

import (
	"time"

	"github.com/bnkamalesh/gonotes/pkg/platform/storage/mongo"
)

// Service defines all the methods implemented by the store
type Service interface {
	// Save saves given data into the store and return the meta info and error if any
	Save(bucket string, data interface{}) (*DocMeta, error)

	// FindOne finds a single document with the provided ID, and loads in result if result is
	// provided
	FindOne(bucket, id string, result interface{}) (map[string]interface{}, error)

	// UpdateOne updates a record with the given ID and new data
	UpdateOne(bucket, id string, data interface{}) (*DocMeta, error)

	// DeleteOne deletes a record with the given ID
	DeleteOne(bucket, id string) (*DocMeta, error)
}

// handlerServices interface defines all the methods required to be a storage service
type handlerServices interface {
	// InsertInfo inserts a new record and return the inserted record's ID
	InsertInfo(bucket string, data interface{}) (string, error)

	// FindOne finds a single document with the provided ID, and loads in result if result is
	// provided
	FindOne(bucket, id string, result interface{}) (map[string]interface{}, error)

	// UpdateOne updates a record with the given ID and new data
	UpdateOne(bucket, id string, data interface{}) error

	// DeleteOne deletes a record with the given ID
	DeleteOne(bucket, id string) error
}

// Config struct holds all the configurations required for the store
type Config struct {
	Name        string
	Username    string
	Password    string
	Hosts       []string
	DialTimeout time.Duration
	Timeout     time.Duration
}

// DocMeta is the document meta generated after performing any store actions
type DocMeta struct {
	Count int
	ID    string
}

// Store holds all the dependencies
type Store struct {
	handler handlerServices
}

// Save saves data into the primary store
func (s *Store) Save(bucket string, data interface{}) (*DocMeta, error) {
	if data == nil {
		return nil, nil
	}
	id, err := s.handler.InsertInfo(bucket, data)
	if err != nil {
		return nil, err
	}

	return &DocMeta{
		ID:    id,
		Count: 1,
	}, nil
}

// FindOne finds a single record based on the provided ID
func (s *Store) FindOne(bucket, id string, result interface{}) (map[string]interface{}, error) {
	return s.handler.FindOne(bucket, id, result)
}

// UpdateOne updates a record with the given ID and new data
func (s *Store) UpdateOne(bucket, id string, data interface{}) (*DocMeta, error) {
	err := s.handler.UpdateOne(bucket, id, data)
	if err != nil {
		return nil, err
	}
	return &DocMeta{ID: id, Count: 1}, nil
}

// DeleteOne deletes a record with the given ID
func (s *Store) DeleteOne(bucket, id string) (*DocMeta, error) {
	err := s.handler.DeleteOne(bucket, id)
	if err != nil {
		return nil, err
	}
	return &DocMeta{ID: id, Count: 1}, nil
}

// New returns a new Service instance
func New(c Config) (Service, error) {
	mongoHandler, err := mongo.New(mongo.Config{
		Name:     c.Name,
		Host:     c.Hosts,
		Username: c.Username,
		Password: c.Password,
		Timeout:  c.Timeout,
	})

	if err != nil {
		return nil, err
	}

	return &Store{
		handler: mongoHandler,
	}, nil
}
