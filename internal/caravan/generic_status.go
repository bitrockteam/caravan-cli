package caravan

type GenericStatus struct {
	GenericProvider
}

func (g GenericStatus) Status() error {
	panic("implement me")
}
