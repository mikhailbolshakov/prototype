package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sherifabdlnaby/configuro"
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

func EnrichWithEnv(filename string, config *Config) error {

	var opts = []configuro.ConfigOptions{ configuro.WithLoadFromEnvVars("CONFIG"), configuro.WithoutLoadFromConfigFile() }

	if _, err := os.Stat(filename); err == nil {
		opts = append(opts, configuro.WithLoadDotEnv(filename))
	} else {
		if os.IsNotExist(err) {
			opts = append(opts, configuro.WithoutLoadDotEnv())
		} else {
			return err
		}
	}

	Loader, err := configuro.NewConfig(opts...)
	if err != nil {
		return err
	}
	return Loader.Load(config)
}
