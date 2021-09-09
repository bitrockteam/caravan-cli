package caravan

type GenericDeploy struct {
	GenericProvider
}

func (g GenericDeploy) Deploy(layer DeployLayer) error {
	panic("implement me")
}
