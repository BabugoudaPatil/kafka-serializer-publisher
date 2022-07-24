package config

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
)

type loader struct {
	overwrite bool
	vipe      *viper.Viper
}

// Load - Builds the configuration from external files
func Load(overwrite bool, files ...string) *Config {
	vipe := viper.New()
	vipe.AutomaticEnv()

	loader := loader{vipe: vipe, overwrite: overwrite}
	loader.loadFiles(files...)

	var cfg Config
	if err := vipe.Unmarshal(&cfg); err != nil { // Handle errors reading the loader file
		panic(fmt.Errorf("Fatal error loader file: %w \n", err))
	}
	return &cfg
}

// LoadFiles - Loads a bunch of loader files into loader
func (c loader) loadFiles(files ...string) {
	for _, file := range files {
		if _, err := os.Stat(file); err == nil {
			c.loadFile(file)
		} else {
			log.Warn().Msgf("Could not load loader file: %s", file)
		}
	}
}

func (c loader) loadFile(file string) {
	c.vipe.SetConfigFile(file)
	var err error
	if c.overwrite {
		err = c.vipe.ReadInConfig()
	} else {
		err = c.vipe.MergeInConfig()
	}
	if err != nil { // Handle errors reading the loader file
		panic(fmt.Errorf("Fatal error loader file: %w \n", err))
	}
}
