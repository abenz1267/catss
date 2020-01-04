package configuration

import (
	"bytes"
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
	File          string `json:"-"`
	content       []byte
	Root          string   `json:"root"`
	Minify        bool     `json:"-"`
	Outputs       []Output `json:"outputs"`
	CreateMissing bool     `json:"-"`
}

func Load(file string) (*Config, error) {
	cfg := &Config{}

	if _, err := os.Stat(file); os.IsNotExist(err) {
		createDummy(file)
	}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return cfg, err
	}

	cfg.content = b

	err = json.Unmarshal(b, cfg)
	if err != nil {
		return cfg, err
	}

	validate(cfg)

	return cfg, err
}

func Update(cfg *Config) (bool, error) {
	b, err := ioutil.ReadFile(cfg.File)
	if err != nil {
		return false, err
	}

	if !bytes.Equal(b, cfg.content) {
		err := json.Unmarshal(b, cfg)
		if err != nil {
			return false, err
		}

		cfg.content = b

		log.Println("Updated config")

		return true, nil
	}

	return false, nil
}

func createDummy(file string) {
	c := `
{
  "root": "",
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

func validate(cfg *Config) {
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
