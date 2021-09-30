package gcp_test

import (
	"caravan-cli/cli"
	"caravan-cli/provider/gcp"
	"context"
	"testing"
)

func TestValidate(t *testing.T) {
	ctx := context.Background()

	type test struct {
		name   string
		error  bool
		desc   string
		region string
	}

	tests := []test{
		{name: "test-me", error: false, desc: "ok", region: "europe-west6"},
		{name: "test", error: true, desc: "name shorter than minimum", region: "europe-west6"},
		{name: "test-me?", error: true, desc: "non supported characters", region: "europe-west6"},
		{name: "-test-me", error: true, desc: "starting with hypen", region: "europe-west6"},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			c, err := cli.NewConfigFromScratch(tc.name, "gcp", tc.region)
			if err != nil {
				t.Fatalf("unable to create config: %s\n", err)
			}
			c.UserEmail = "test.name@test.me"
			_, err = gcp.New(ctx, c)
			if err == nil && tc.error || err != nil && !tc.error {
				t.Errorf("something wen wrong: want %t but got %s", tc.error, err)
			}
		})
	}
}
