package configuration

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// Environment structure
type Env struct {
	DataBase
	Redis
}

// DB configuration
type DataBase struct {
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

// Server configuration structure, Port string
type Redis struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// Init environment with configuration
func InitEnv(p string) (*Env, error) {
	file, err := os.Open(p)

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Println(err)
		}
	}(file)

	var env Env
	if err != nil {
		return &env, err
	}

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&env)
	if err != nil {
		panic(err)
	}

	return &env, nil
}
