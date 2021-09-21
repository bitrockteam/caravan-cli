// Copyright © 2021 Bitrock s.r.l. <devops@bitrock.it>
package cmd

import (
	"caravan/internal/aws"
	"caravan/internal/caravan"
	"caravan/internal/gcp"
	"fmt"
)

func getProvider(c *caravan.Config) (caravan.Provider, error) {
	var p caravan.Provider
	var err error
	switch c.Provider {
	case caravan.AWS:
		p, err = aws.New(c)
	case caravan.GCP:
		p, err = gcp.New(c)
	// case caravan.Azure:
	//	p, err = azure.New(c)
	default:
		p, err = nil, fmt.Errorf("unknown provider")
	}
	if err != nil {
		return nil, err
	}
	return p, nil
}
