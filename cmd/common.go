package cmd

import (
	"caravan-cli/cli"
	"caravan-cli/provider"
	"caravan-cli/provider/aws"
	"caravan-cli/provider/azure"
	"caravan-cli/provider/gcp"
	"context"
	"fmt"
)

func getProvider(ctx context.Context, c *cli.Config) (provider.Provider, error) {
	var p provider.Provider
	var err error
	switch c.Provider {
	case provider.AWS:
		p, err = aws.New(ctx, c)
	case provider.GCP:
		p, err = gcp.New(ctx, c)
	case provider.Azure:
		p, err = azure.New(ctx, c)
	default:
		p, err = nil, fmt.Errorf("unknown provider")
	}
	if err != nil {
		return nil, err
	}
	return p, nil
}
