// Package configs handles all the configurations required for this app
package configs

import (
	"github.com/bnkamalesh/webgo"
)

// Webgo returns the configurations required for webgo
func Webgo() *webgo.Config {
	return &webgo.Config{
		Host: "",
		Port: "8080",
	}
}
