package config

import "os"

// Represents database server and credentials
type Config struct {
	Server   string
	Database string
}

// Get the config info
func (c *Config) Read() {
	c.Server = os.Getenv("dserver")
	c.Database = os.Getenv("ddb")
}
