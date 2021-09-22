// Copyright © 2021 Bitrock s.r.l. <devops@bitrock.it>
package cmd

import (
	"caravan-cli/config"
	"caravan-cli/provider"
	"caravan-cli/provider/aws"
	"caravan-cli/provider/gcp"
	"fmt"
)

func getProvider(c *config.Config) (provider.Provider, error) {
	var p provider.Provider
	var err error
	switch c.Provider {
	case provider.AWS:
		p, err = aws.New(c)
	case provider.GCP:
		p, err = gcp.New(c)
	// case provider.Azure:
	//	p, err = azure.New(c)
	default:
		p, err = nil, fmt.Errorf("unknown provider")
	}
	if err != nil {
		return nil, err
	}
	return p, nil
}
