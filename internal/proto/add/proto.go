package add

import (
	"fmt"
	"os"
	"path"
)

// Proto is a proto generator.
type Proto struct {
	FileName  string
	Path      string
	Service   string
	Package   string
	GoPackage string
}

// Generate generate a proto template.
func (p *Proto) Generate() (err error) {

	var body []byte
	if body, err = p.execute(); err != nil {
		return
	}

	var dir string
	if dir, err = os.Getwd(); err != nil {
		return
	}

	to := path.Join(dir, p.Path)
	if _, err = os.Stat(to); os.IsNotExist(err) {
		if err = os.MkdirAll(to, 0o700); err != nil {
			return
		}
	}

	name := path.Join(to, p.FileName)
	if _, err = os.Stat(name); !os.IsNotExist(err) {
		return fmt.Errorf("%s already exists", p.FileName)
	}
	return os.WriteFile(name, body, 0o644)
}
