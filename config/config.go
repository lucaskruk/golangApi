package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

var CfgPath string = "config/config.yml"

type Config struct {
	Ships []struct { //almacena las posiciones de los satelites y sus nombres
		Name string  `yaml:"name"`
		X    float32 `yaml:"x"`
		Y    float32 `yaml:"y"`
	} `yaml:"ships"`

	Server struct {
		// la ip de la maquina local asociada al server
		Host string `yaml:"host"`

		// el puerto asociado al server en la maquina local
		Port    string `yaml:"port"`
		Timeout struct {
			// Timeout general para desconexiones exitosas
			Server time.Duration `yaml:"server"`

			// Tiempo de espera antes de cancelar la op de escritura en un server HTTP
			Write time.Duration `yaml:"write"`

			// Tiempo de espera antes de cancelar la op de lectura en un server HTTP
			Read time.Duration `yaml:"read"`

			// Tiempo de espera para cerrar sesiones inactivas
			Idle time.Duration `yaml:"idle"`
		} `yaml:"timeout"`
	} `yaml:"server"`
}

// NewConfig devuelve el struct de configuracion cargado en memoria
func NewConfig(configPath string) (*Config, error) {

	config := &Config{}

	// abro el archivo de configuracion
	file, err := os.Open(configPath)
	if err != nil { //arreglo para los tests
		file, err = os.Open("../" + configPath)
		if err != nil {
			return nil, err
		}
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	// decodifico el YAML del archivo
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

// ValidateConfigPath asegura, que la ruta es un archivo legible
func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' es un directorio, no un archivo", path)
	}
	return nil
}
