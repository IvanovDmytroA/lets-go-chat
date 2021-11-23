package configuration

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

const configPath string = "../configs/config.yml"

// Environment structure
type Env struct {
	DataBase
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

// Init environment with configuration
func InitEnv() (*Env, error) {
	file, err := os.Open(configPath)

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Println(err)
		}
	}(file)

	if err != nil {
		panic(err)
	}

	decoder := yaml.NewDecoder(file)
	var env Env
	err = decoder.Decode(&env)
	if err != nil {
		panic(err)
	}

	return &env, nil
}
