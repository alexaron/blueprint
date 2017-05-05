// Package env reads the application settings.
package env

import (
	"encoding/json"

	"fmt"
	"github.com/blue-jay/core/asset"
	"github.com/blue-jay/core/email"
	"github.com/blue-jay/core/form"
	"github.com/blue-jay/core/generate"
	"github.com/blue-jay/core/jsonconfig"
	"github.com/blue-jay/core/server"
	"github.com/blue-jay/core/session"
	"github.com/blue-jay/core/view"
)

// *****************************************************************************
// Application Settings
// *****************************************************************************

// Info structures the application settings.
type Info struct {
	Asset      asset.Info    `json:"Asset"`
	Email      email.Info    `json:"Email"`
	Form       form.Info     `json:"Form"`
	Generation generate.Info `json:"Generation"`
	PostgreSQL PostgreSQL    `json:"PostgreSQL"`
	Server     server.Info   `json:"Server"`
	Session    session.Info  `json:"Session"`
	Template   view.Template `json:"Template"`
	View       view.Info     `json:"View"`
	path       string
}

// PostgreSQL stores postgres configuration.
type PostgreSQL struct {
	Hostname string `json:"Hostname"`
	Port     int    `json:"Port"`
	Username string `json:"Username"`
	Password string `json:"Password"`
	Database string `json:"Database"`
	SSLMode  bool   `json:"SSLMode"`
}

// DSN creates a connection string
func (p PostgreSQL) DSN() string {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		p.Username,
		p.Password,
		p.Hostname,
		p.Port,
		p.Database,
	)
	if !p.SSLMode {
		dsn += "?sslmode=disable"
	}
	return dsn
}

// Path returns the env.json path
func (c *Info) Path() string {
	return c.path
}

// ParseJSON unmarshals bytes to structs
func (c *Info) ParseJSON(b []byte) error {
	return json.Unmarshal(b, &c)
}

// New returns a instance of the application settings.
func New(path string) *Info {
	return &Info{
		path: path,
	}
}

// LoadConfig reads the configuration file.
func LoadConfig(configFile string) (*Info, error) {
	// Create a new configuration with the path to the file
	config := New(configFile)

	// Load the configuration file
	err := jsonconfig.Load(configFile, config)

	// Return the configuration
	return config, err
}
