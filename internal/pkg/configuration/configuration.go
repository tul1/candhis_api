package configuration

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

func loadConfigFromYML[T any](configFilePath string, config T) error {
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		return err
	}

	validate := validator.New()
	return validate.Struct(config)
}

func Load[T any](configFile string) (*T, error) {
	config := new(T)
	if err := loadConfigFromYML(configFile, config); err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}
	return config, nil
}
