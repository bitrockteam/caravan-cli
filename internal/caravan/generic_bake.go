package caravan

type GenericBake struct {
	GenericProvider
}

func (g GenericBake) Bake() error {
	panic("implement me")
}
