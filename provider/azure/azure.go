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
	baking := caravan.Template{
		Name: "baking-vars",
		Text: bakingTfVarsTmpl,
		Path: a.Caravan.WorkdirBakingVars,
	}
	infra := caravan.Template{
		Name: "infra-vars",
		Text: infraTfVarsTmpl,
		Path: a.Caravan.WorkdirInfraVars,
	}
	infraBackend := caravan.Template{
		Name: "infra-backend",
		Text: infraBackendTmpl,
		Path: a.Caravan.WorkdirInfraBackend,
	}
	platform := caravan.Template{
		Name: "platform-vars",
		Text: platformTfVarsTmpl,
		Path: a.Caravan.WorkdirPlatformVars,
	}
	platformBackend := caravan.Template{
		Name: "platform-backend",
		Text: platformBackendTmpl,
		Path: a.Caravan.WorkdirPlatformBackend,
	}
	applicationSupport := caravan.Template{
		Name: "application-vars",
		Text: applicationTfVarsTmpl,
		Path: a.Caravan.WorkdirApplicationVars,
	}
	applicationSupportBackend := caravan.Template{
		Name: "application-backend",
		Text: applicationSupportBackendTmpl,
		Path: a.Caravan.WorkdirApplicationBackend,
	}

	return []caravan.Template{
		baking,
		infra,
		infraBackend,
		platform,
		platformBackend,
		applicationSupport,
		applicationSupportBackend,
	}, nil
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
