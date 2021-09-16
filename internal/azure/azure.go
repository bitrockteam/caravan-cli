package azure

import "caravan/internal/caravan"

type Azure struct {
	caravan.GenericProvider
	caravan.GenericBake
	caravan.GenericDeploy
	caravan.GenericStatus
}

func (a Azure) GetTemplates() ([]caravan.Template, error) {
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
// func (a Azure) Deploy(layer caravan.DeployLayer) error {
//	panic("implement me")
// }

func (a Azure) Clean() error {
	panic("implement me")
}

// func (a Azure) Status() error {
//	panic("implement me")
// }
