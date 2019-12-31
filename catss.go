package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/abenz1267/catss/pkg/cat"
	"github.com/abenz1267/catss/pkg/configuration"
	"github.com/abenz1267/catss/pkg/filewatcher"
	"github.com/pborman/getopt/v2"
)

var (
	cfgFile = "catss.json"
	watch   bool
	minify  bool
)

func init() {
	getopt.Flag(&cfgFile, 'c', "config file")
	getopt.Flag(&watch, 'w', "watch for changes")
	getopt.Flag(&minify, 'm', "minify css")
}

func main() {
	getopt.Parse()

	wDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := configuration.Load(filepath.Join(wDir, cfgFile))
	if err != nil {
		log.Fatal(err)
	}

	cfg.Minify = minify
	cfg.Root = filepath.Join(wDir, cfg.Root)

	err = cat.Load(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if watch {
		filewatcher.Watch(cfg)
	}
}
