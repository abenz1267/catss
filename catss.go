package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/abenz1267/catss/pkg/cat"
	"github.com/abenz1267/catss/pkg/configuration"
	"github.com/abenz1267/catss/pkg/filewatcher"
)

func main() {
	cfgFile := flag.String("cfg", "catss.json", "Config file. Root is the working dir.")

	flag.Parse()

	wDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := configuration.Load(filepath.Join(wDir, *cfgFile))
	if err != nil {
		log.Fatal(err)
	}

	cfg.Root = filepath.Join(wDir, cfg.Root)

	err = cat.Load(cfg)
	if err != nil {
		log.Fatal(err)
	}

	filewatcher.Watch(cfg)
}
