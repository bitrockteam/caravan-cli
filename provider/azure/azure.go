package azure

import (
	"caravan-cli/cli"
	"caravan-cli/provider"
)

type Azure struct {
	provider.GenericProvider
}

func (a Azure) GetTemplates() ([]cli.Template, error) {
	panic("implement me")
}

func (a Azure) ValidateConfiguration() error {
	panic("implement me")
}

func (a Azure) InitProvider() error {
	panic("implement me")
}

// func (a Azure) Bake() error {
//	panic("implement me")
// }
//
// func (a Azure) Deploy(layer cli.DeployLayer) error {
//	panic("implement me")
// }

func (a Azure) Clean() error {
	panic("implement me")
}

// func (a Azure) Status() error {
//	panic("implement me")
// }
