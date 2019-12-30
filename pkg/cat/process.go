package cat

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/abenz1267/catss/pkg/configuration"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
)

const PERM = 777
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

func Load(cfg configuration.Config) error {
	var err error

	for _, v := range cfg.Outputs {
		no := output{file: v.File}

		for _, m := range v.Files {
			f := file{name: m}

			f.content, err = ioutil.ReadFile(filepath.Join(cfg.Root, m+EXT))
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

	if min {
		minified, err := minifier.Bytes(MIME, b.Bytes())
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(out.file+EXT, minified, PERM)
		if err != nil {
			return err
		}

		return nil
	}

	fn := out.file + EXT

	if _, err := os.Stat(fn); !os.IsNotExist(err) {
		err := os.Remove(fn)
		if err != nil {
			return err
		}
	}

	err = ioutil.WriteFile(fn, b.Bytes(), PERM)
	if err != nil {
		return err
	}

	log.Printf("Created file: %s", fn)

	return nil
}

func Update(file string, root string, minify bool) error {
	cleaned := strings.Replace(file, root, "", 1)
	cleaned = strings.Replace(cleaned, EXT, "", 1)
	cleaned = strings.TrimPrefix(cleaned, "/")

	for i, v := range o {
		for n, m := range v.files {
			if cleaned == m.name {
				b, err := ioutil.ReadFile(file)
				if err != nil {
					return err
				}

				if !isEqual(b, m.content) {
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

func isEqual(a, b []byte) bool {
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
