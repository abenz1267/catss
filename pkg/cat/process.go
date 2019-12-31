package cat

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/abenz1267/catss/pkg/configuration"
	"github.com/abenz1267/catss/pkg/util"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
)

const EXT = ".css"
const MIME = "test/css"

type output struct {
	file  string
	files []file
}

type file struct {
	name    string
	content []byte
}

var o []output

func Load(cfg *configuration.Config) error {
	var err error
	o = nil

	for _, v := range cfg.Outputs {
		no := output{file: v.File}

		for _, m := range v.Files {
			f := file{name: m}

			path := filepath.Join(cfg.Root, m+EXT)

			if cfg.CreateMissing {
				if _, err := os.Stat(path); os.IsNotExist(err) {
					file, err := os.Create(path)
					if err != nil {
						return err
					}
					file.Close()

					log.Printf("Created file: %s", path)
				}
			}

			f.content, err = ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			no.files = append(no.files, f)
		}

		o = append(o, no)
	}

	return cat(cfg.Minify)
}

func cat(minify bool) error {
	for _, v := range o {
		err := writeOutput(v, minify)
		if err != nil {
			return err
		}
	}

	return nil
}

func writeOutput(out output, min bool) error {
	var b bytes.Buffer
	var err error

	minifier := minify.New()
	minifier.AddFunc(MIME, css.Minify)

	for _, m := range out.files {
		b.Write(m.content)
	}

	fn := out.file + EXT
	f, err := os.Create(fn)
	if err != nil {
		return err
	}

	defer f.Close()

	if min {
		minified, err := minifier.Bytes(MIME, b.Bytes())
		if err != nil {
			return err
		}

		_, err = f.Write(minified)
		if err != nil {
			return err
		}
	} else {
		_, err = f.Write(b.Bytes())
		if err != nil {
			return err
		}
	}

	log.Printf("Created file: %s", fn)

	return nil
}

func Update(file string, root string, minify bool) error {
	sep := "/"

	if runtime.GOOS == "windows" {
		sep = "\\"
	}

	cleaned := strings.Replace(file, root, "", 1)
	cleaned = strings.Replace(cleaned, EXT, "", 1)
	cleaned = strings.TrimPrefix(cleaned, sep)

	for i, v := range o {
		for n, m := range v.files {
			if cleaned == m.name {
				b, err := ioutil.ReadFile(file)
				if err != nil {
					return err
				}

				if !util.IsEqual(b, m.content) {
					log.Printf("Updating file: %s", file)
					o[i].files[n].content = b

					err := writeOutput(o[i], minify)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}
