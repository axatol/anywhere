package config

import (
	_ "embed"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	Values config
)

func resolveRawConfig() []byte {
	var locationsToTry []string
	if loc, ok := os.LookupEnv("ANYWHERE_CONFIG_FILE"); ok {
		locationsToTry = []string{loc}
	} else {
		locationsToTry = []string{
			"./config.yml",
			"./example.config.yml",
		}
	}

	for _, loc := range locationsToTry {
		if _, err := os.Stat(loc); err != nil {
			continue
		}

		rawConfig, err := os.ReadFile(loc)
		if err != nil {
			panic(err)
		}

		return rawConfig
	}

	return nil
}

func initConfig() {
	if rawConfig := resolveRawConfig(); rawConfig != nil {
		if err := yaml.Unmarshal(rawConfig, &Values); err != nil {
			panic("failed to parse config file")
		}
	}

	// applying defaults to missing values
	Values.ApplyDefaults()

	// overriding with environment variables
	Values.ApplyFromEnvironment()

	// overriding with flags
	Values.ApplyFromFlags()
}
