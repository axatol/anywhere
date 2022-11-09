package config

import (
	_ "embed"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	Values config
)

func resolveRawConfig() []byte {
	var locationsToTry []string

	if loc := envStr("CONFIG_FILEPATH"); loc != nil {
		locationsToTry = []string{*loc}
	} else {
		locationsToTry = []string{
			"./config.local.yml",
			"./config.yml",
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

		fmt.Printf("INFO: using config at: %s\n", loc)
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
