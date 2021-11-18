package configuration

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

const configPath string = "configs/config.yml"

// Environment structure
type Env struct {
	DataBase
}

// DB configuration
type DataBase struct {
	Type     string `yaml:"database.type"`
	Host     string `yaml:"database.host"`
	Port     int    `yaml:"database.port"`
	Name     string `yaml:"database.name"`
	User     string `yaml:"database.user"`
	Password string `yaml:"database.password"`
}

// Init environment with configuration
func InitEnv() (*Env, error) {
	file, err := getConfigFile(configPath)
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

func getConfigFile(filePath string) (*os.File, error) {
	file, err := os.Open(configPath)

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Println(err)
		}
	}(file)

	return file, err
}
