// Package storage handles the primary data store
package storage

import (
	"errors"
	"time"

	"github.com/bnkamalesh/notes/pkg/platform/storage/mongo"
)

var (
	// ErrNotFound is returned if the record was not found in the storage
	ErrNotFound = errors.New("Record not found")
)

// Service defines all the methods implemented by the store
type Service interface {
	// Save saves given data into the store and return the meta info and error if any
	Save(bucket string, data interface{}) (*DocMeta, error)

	// Update updates the first record matching the given query and new data
	Update(bucket string, query interface{}, data interface{}) error

	// Delete deletes the first record matching the provided query
	Delete(collectionName string, query interface{}) error

	// Find finds all the records matching the query
	Find(bucket string, query, selectFields interface{}, sort []string, start, limit int, result interface{}) ([]map[string]interface{}, error)

	// FindOne finds the first matching document for the given query
	FindOne(bucket string, query, selectFields interface{}, sort []string, result interface{}) (map[string]interface{}, error)
}

// handlerServices interface defines all the methods required to be a storage service
type handlerServices interface {
	// InsertInfo inserts a new record and return the inserted record's ID
	InsertInfo(bucket string, data interface{}) (string, error)

	// Update updates the first record matching the given query and new data
	Update(bucket string, query interface{}, data interface{}) error

	// Delete delets the first record matching the query
	Delete(collectionName string, query interface{}) error

	// Find finds all the records matching the query
	Find(bucket string, query, selectFields interface{}, sort []string, start, limit int, result interface{}) ([]map[string]interface{}, error)

	// FindOne finds the first matching document for the given query
	FindOne(bucket string, query, selectFields interface{}, sort []string, result interface{}) (map[string]interface{}, error)
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

// Find finds all the records based on the provided query
func (s *Store) Find(bucket string, query, selectFields interface{}, sort []string, start, limit int, result interface{}) ([]map[string]interface{}, error) {
	out, err := s.handler.Find(bucket, query, selectFields, sort, start, limit, result)
	if err == mongo.ErrNotFound {
		return nil, ErrNotFound
	}
	return out, err
}

// FindOne finds the first document matching the provided query
func (s *Store) FindOne(bucket string, query, selectFields interface{}, sort []string, result interface{}) (map[string]interface{}, error) {
	out, err := s.handler.FindOne(bucket, query, selectFields, sort, result)
	if err == mongo.ErrNotFound {
		return nil, ErrNotFound
	}
	return out, err
}

// Update updates the first record matching the query
func (s *Store) Update(bucket string, query interface{}, data interface{}) error {
	err := s.handler.Update(bucket, query, data)
	if err != nil {
		if err == mongo.ErrNotFound {
			return ErrNotFound
		}
		return err
	}
	return nil
}

// Delete deletes the first record matching the query
func (s *Store) Delete(bucket string, query interface{}) error {
	err := s.handler.Delete(bucket, query)
	if err != nil {
		if err == mongo.ErrNotFound {
			return ErrNotFound
		}
		return err
	}
	return nil
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
