// Package configs handles all the configurations required for this app
package configs

import (
	"os"
	"time"

	"github.com/bnkamalesh/webgo"

	"github.com/bnkamalesh/notes/pkg/platform/cache"
	"github.com/bnkamalesh/notes/pkg/platform/storage"
)

// Logs returns which all logs should be enabled
func Logs() []string {
	return []string{"all"}
}

// Webgo returns the configurations required for webgo
func Webgo() *webgo.Config {
	return &webgo.Config{
		Host: "",
		Port: os.Getenv("notes_app_httpPort"),
	}
}

// Store returns the configuration required for the primary datastore
func Store() storage.Config {
	return storage.Config{
		Name:                 os.Getenv("notes_db_name"),
		Hosts:                []string{os.Getenv("notes_db_host")},
		Username:             os.Getenv("notes_user"),
		Password:             os.Getenv("notes_password"),
		AuthenticationSource: os.Getenv("notes_db_authsource"),
		Timeout:              time.Second * 3,
		DialTimeout:          time.Second * 15,
	}
}

// Cache returns the configuration required for cache
func Cache() cache.Config {
	return cache.Config{
		Name:         "0",
		Hosts:        []string{os.Getenv("notes_cache_host")},
		DialTimeout:  time.Second * 15,
		ReadTimeout:  time.Millisecond * 25,
		WriteTimeout: time.Millisecond * 75,
	}
}
