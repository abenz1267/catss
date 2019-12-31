package configuration

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
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

	if _, err := os.Stat(file); os.IsNotExist(err) {
		createDummy(file)
	}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return cfg, err
	}

	json.Unmarshal(b, &cfg)

	validate(cfg)

	return cfg, err
}

func createDummy(file string) {
	c := `
{
  "root": "",
  "minify": false,
  "outputs": [
    {
      "file": "style",
      "files": ["first", "second"]
    }
  ]
}
	`

	f, err := os.Create(file)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err = f.WriteString(strings.TrimSpace(c))
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal("Config file created. You need to edit it before using Catss.")
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
