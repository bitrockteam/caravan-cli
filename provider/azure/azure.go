// Microsoft Azure provider.
package azure

import (
	"caravan-cli/cli"
	"caravan-cli/provider"
)

type Azure struct {
	provider.GenericProvider
}

func New(c *caravan.Config) (Azure, error) {
	a := Azure{}
	a.Caravan = c
	if err := a.ValidateConfiguration(); err != nil {
		return a, err
	}
	return a, nil
}

func (a Azure) GetTemplates() ([]cli.Template, error) {
	panic("implement me")
}

func (a Azure) ValidateConfiguration() error {
	return nil
}

func (a Azure) InitProvider() error {
	panic("implement me")
}

func (a Azure) Bake() error {
	panic("implement me")
}

func (a Azure) Deploy(layer caravan.DeployLayer) error {
	panic("implement me")
}

func (a Azure) Destroy(layer caravan.DeployLayer) error {
	panic("implement me")
}

func (a Azure) CleanProvider() error {
	panic("implement me")
}

func (a Azure) Status() error {
	panic("implement me")
}
