package filewatcher

import (
	"log"
	"path/filepath"

	"github.com/abenz1267/catss/pkg/cat"
	"github.com/abenz1267/catss/pkg/configuration"
	"github.com/fsnotify/fsnotify"
)

func Watch(cfg configuration.Config) {
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

				cat.Update(event.Name, cfg.Root, cfg.Minify)

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

	for _, v := range cfg.Outputs {
		for _, m := range v.Files {
			err = watcher.Add(filepath.Join(cfg.Root, m+cat.EXT))
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	<-done
}
