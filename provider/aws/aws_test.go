package aws_test

import (
	"caravan-cli/cli"
	"caravan-cli/provider/aws"
	"testing"
)

func TestValidate(t *testing.T) {
	type test struct {
		name  string
		error bool
		desc  string
	}

	tests := []test{
		{name: "test-me", error: false, desc: "ok"},
		{name: "test567890123", error: true, desc: "name too long"},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			c, err := cli.NewConfigFromScratch(tc.name, "aws", "eu-south-1")
			if err != nil {
				t.Fatalf("unable to create config: %s\n", err)
			}
			_, err = aws.New(c)
			if err == nil && tc.error || err != nil && !tc.error {
				t.Errorf("something wen wrong: want %t but got %s", tc.error, err)
			}
		})
	}
}
