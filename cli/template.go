package cli

import (
	caravan "caravan-cli/config"
	"os"
	"path/filepath"
	"text/template"
)

type Template struct {
	Name string
	Text string
	Path string
}

func (t Template) Render(c *caravan.Config) error {
	temp, err := template.New(t.Name).Parse(t.Text)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(t.Path), 0777); err != nil {
		return err
	}
	f, err := os.Create(t.Path)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := temp.Execute(f, c); err != nil {
		return err
	}
	return nil
}
