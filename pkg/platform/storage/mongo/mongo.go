package mongo

import (
	"errors"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/google/uuid"
)

var (
	randomValue = uuid.New().String()
	randomKey   = uuid.New().String()
	randomQuery = bson.M{randomKey: randomValue}
	// ErrNotFound is returned when the document was not found in Mongo collection
	ErrNotFound = errors.New("Document not found")
)

// Config holds the config required for MongoDB
type Config struct {
	AppName          string   `json:"appName,omitempty"`
	Name             string   `json:"name,omitempty"`
	Host             []string `json:"host,omitempty"`
	Port             string   `json:"port,omitempty"`
	ConnectionString string   `json:"connectionString,omitempty"`
	AuthSource       string   `json:"authSource,omitempty"`
	ReplicaSet       string   `json:"replicaSet,omitempty"`

	Username       string             `json:"username,omitempty"`
	Password       string             `json:"password,omitempty"`
	Direct         bool               `json:"direct,omitempty"`
	Timeout        time.Duration      `json:"timeout,omitempty"`
	ReadPreference mgo.ReadPreference `json:"readPreference,omitempty"`
}

// Handler DB Sessions are maintained inside a struct for better caching of the data stores
// developed based on the stackoverflow answer:
// http://stackoverflow.com/questions/26574594/best-practice-to-maintain-a-mgo-session
type Handler struct {
	DBName  string
	session *mgo.Session
}

// Clone the master session and return
func (ms *Handler) getSession() *mgo.Session {
	return ms.session.Copy()
}

//SessionCollection gets the appropriate MongoDB collection
func (ms *Handler) sessionCollection(collection string) (*mgo.Session, *mgo.Collection) {
	s := ms.getSession()
	c := s.DB(ms.DBName).C(collection)
	return s, c
}

// New returns a new MongoDB handler instance with all the configurations set
func New(c Config) (*Handler, error) {
	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Database:       c.Name,
		ReplicaSetName: c.ReplicaSet,
		Addrs:          c.Host,
		Username:       c.Username,
		Password:       c.Password,
		Timeout:        c.Timeout,
		Source:         c.AuthSource,
		Direct:         c.Direct,
		AppName:        c.AppName,
		ReadPreference: &c.ReadPreference,
	})

	if err != nil {
		return nil, err
	}

	session.SetSafe(&mgo.Safe{WMode: "majority"})
	err = session.Ping()
	if err != nil {
		return nil, err
	}
	return &Handler{DBName: c.Name, session: session}, nil
}

// InsertInfo inserts a new document and return inserted document's ID
func (ms *Handler) InsertInfo(collectionName string, data interface{}) (string, error) {
	session, collection := ms.sessionCollection(collectionName)
	defer session.Close()

	// randomKey,randomValue pair is used to ensure that a new document is inserted every
	// time Upsert is called. This workaround is done to get the inserted document's ID.
	// Upsert is the only method available which returns the inserted document's ID
	info, err := collection.Upsert(randomQuery, data)
	if err != nil {
		return "", err
	}

	id := ""
	if objID, ok := info.UpsertedId.(bson.ObjectId); ok {
		id = objID.Hex()
	} else {
		return "", errors.New("Invalid ID received")
	}

	return id, nil
}

// FindOne finds a document with the given ID
func (ms *Handler) FindOne(collectionName, id string, result interface{}) (map[string]interface{}, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, errors.New("Invalid ID provided")
	}
	objID := bson.ObjectIdHex(id)

	session, collection := ms.sessionCollection(collectionName)
	defer session.Close()

	if result != nil {
		return nil, collection.FindId(objID).One(result)
	}

	out := make(map[string]interface{}, 0)
	err := collection.FindId(objID).One(&out)
	if err == mgo.ErrNotFound {
		return nil, ErrNotFound
	}
	return out, err
}

// UpdateOne updates a document with the given ID
func (ms *Handler) UpdateOne(collectionName, id string, data interface{}) error {
	if !bson.IsObjectIdHex(id) {
		return errors.New("Invalid ID provided")
	}
	objID := bson.ObjectIdHex(id)

	session, collection := ms.sessionCollection(collectionName)
	defer session.Close()

	return collection.UpdateId(objID, data)
}

// DeleteOne deletes a document with the given ID
func (ms *Handler) DeleteOne(collectionName, id string) error {
	if !bson.IsObjectIdHex(id) {
		return errors.New("Invalid ID provided")
	}
	objID := bson.ObjectIdHex(id)

	session, collection := ms.sessionCollection(collectionName)
	defer session.Close()

	return collection.RemoveId(objID)
}
