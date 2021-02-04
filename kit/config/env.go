package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func Env(filename string) (map[string]string, error) {

	var envMap map[string]string
	if _, err := os.Stat(filename); err == nil || !os.IsNotExist(err) {
		envMap, err = godotenv.Read(filename)
		if err != nil {
			return nil, fmt.Errorf("error loading .env: %s", err.Error())
		}
	}
	return envMap, nil
}
