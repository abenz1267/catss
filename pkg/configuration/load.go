package configuration

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
)

const msg = "Paths can not start with a '/'"

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

	validate(cfg)

	return cfg, err
}

func validate(cfg Config) {
	cp(cfg.Root)

	for _, v := range cfg.Outputs {
		cp(v.File)

		for _, m := range v.Files {
			cp(m)
		}
	}

}

func cp(v string) {
	if strings.HasPrefix(v, "/") {
		log.Fatal(msg)
	}
}
