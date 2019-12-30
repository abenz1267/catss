package configuration

import (
	"encoding/json"
	"io/ioutil"
)

type Output struct {
	File  string   `json:"file"`
	Files []string `json:"files"`
}

type Config struct {
	Root    string   `json:"root"`
	Minify  bool     `json:"minify"`
	Outputs []Output `json:"outputs"`
}

func Load(file string) (Config, error) {
	var cfg Config
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return cfg, err
	}

	json.Unmarshal(b, &cfg)

	return cfg, err
}
