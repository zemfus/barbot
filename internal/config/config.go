package config

import (
	"errors"
	"gopkg.in/yaml.v3"
	"os"
)

func New(pathConfig string) (*Configuration, error) {
	file, err := readFile(pathConfig)
	if err != nil {
		return nil, err
	}

	config := new(Configuration)
	if err = yaml.Unmarshal(file, config); err != nil {
		return nil, errors.Join(errors.New("Unmarshal file error: "), err)
	}

	return config, nil
}

func readFile(pathConfig string) ([]byte, error) {
	file, err := os.ReadFile(pathConfig)
	if err != nil {
		return nil, errors.Join(errors.New("Read file error: "), err)
	}

	return file, nil
}
