// Package configs handles all the configurations required for this app
package configs

import (
	"time"

	"github.com/bnkamalesh/webgo"

	"github.com/bnkamalesh/gonotes/pkg/platform/storage"
)

// Webgo returns the configurations required for webgo
func Webgo() *webgo.Config {
	return &webgo.Config{
		Host: "",
		Port: "8080",
	}
}

// Store returns the configuration required for the primary datastore
func Store() storage.Config {
	return storage.Config{
		Name:        "gonotes",
		Hosts:       []string{"127.0.0.1:27017"},
		Timeout:     time.Second * 3,
		DialTimeout: time.Second * 15,
	}
}
