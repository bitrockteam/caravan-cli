package checker

import (
	"context"
)

type NomadChecker struct {
	GenericChecker
}

func NewNomadChecker(u, ca string, options ...func(*GenericChecker)) (nc NomadChecker, err error) {
	client, err := TLSClient(ca)
	if err != nil {
		return nc, err
	}
	gc := NewGenericChecker(u+"/v1/status/leader", client)
	for _, op := range options {
		if op != nil {
			op(gc)
		}
	}
	return NomadChecker{
		GenericChecker: *gc,
	}, nil
}

func (n NomadChecker) Status(ctx context.Context) bool {
	u := "/v1/status/leader"
	return n.GenericChecker.CheckURL(ctx, u)
}

func (n NomadChecker) Version(ctx context.Context) string {
	// TODO find endpoint
	return "missing endpoint"
}
