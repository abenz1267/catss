package filewatcher

import (
	"log"
	"path/filepath"

	"github.com/abenz1267/catss/pkg/cat"
	"github.com/abenz1267/catss/pkg/configuration"
	"github.com/fsnotify/fsnotify"
)

func Watch(cfg *configuration.Config) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Name != cfg.File {
					cat.Update(event.Name, cfg.Root, cfg.Minify)
				} else {
					updated, err := configuration.Update(cfg)
					if err != nil {
						log.Fatal(err)
					}

					if updated {
						err = cat.Load(cfg)
						if err != nil {
							log.Fatal(err)
						}

						addFiles(watcher, cfg)
					}
				}

				err = watcher.Add(event.Name)
				if err != nil {
					log.Fatal(err)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	addFiles(watcher, cfg)

	err = watcher.Add(cfg.File)
	if err != nil {
		log.Fatal(err)
	}

	<-done
}

func addFiles(watcher *fsnotify.Watcher, cfg *configuration.Config) {
	for _, v := range cfg.Outputs {
		for _, m := range v.Files {
			err := watcher.Add(filepath.Join(cfg.Root, m+cat.EXT))
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
