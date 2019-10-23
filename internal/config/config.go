package config

import (
	"os"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

type Conf struct {
	FwName      string `yaml:"firmware"`
	NetworkType string `yaml:"net"`
	NetworkPort string `yaml:"port"`
}

func NewConfig(filename string) Conf {

	conf := &Conf{}

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal().Err(err).Caller().Msg("")
	}

	if err = yaml.NewDecoder(f).Decode(conf); err != nil {
		log.Fatal().Err(err).Caller().Msg("")
	}

	return *conf
}
