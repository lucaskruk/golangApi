package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

var CfgPath string = "config/config.yml"

type Config struct {
	RebelShip1 struct {
		Name string  `yaml:"name"`
		X    float32 `yaml:"x"`
		Y    float32 `yaml:"y"`
	} `yaml:"rebelship1"`

	RebelShip2 struct {
		Name string  `yaml:"name"`
		X    float32 `yaml:"x"`
		Y    float32 `yaml:"y"`
	} `yaml:"rebelship2"`

	RebelShip3 struct {
		Name string  `yaml:"name"`
		X    float32 `yaml:"x"`
		Y    float32 `yaml:"y"`
	} `yaml:"rebelship3"`

	Server struct {
		// Host is the local machine IP Address to bind the HTTP Server to
		Host string `yaml:"host"`

		// Port is the local machine TCP Port to bind the HTTP Server to
		Port    string `yaml:"port"`
		Timeout struct {
			// Server is the general server timeout to use
			// for graceful shutdowns
			Server time.Duration `yaml:"server"`

			// Write is the amount of time to wait until an HTTP server
			// write opperation is cancelled
			Write time.Duration `yaml:"write"`

			// Read is the amount of time to wait until an HTTP server
			// read operation is cancelled
			Read time.Duration `yaml:"read"`

			// Read is the amount of time to wait
			// until an IDLE HTTP session is closed
			Idle time.Duration `yaml:"idle"`
		} `yaml:"timeout"`
	} `yaml:"server"`
}

// NewConfig returns a new decoded Config struct
func NewConfig(configPath string) (*Config, error) {
	// Create config structure
	config := &Config{}

	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

// ValidateConfigPath just makes sure, that the path provided is a file,
// that can be read
func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

func init() {

}
