package gcp_test

import (
	"caravan/internal/caravan"
	"caravan/internal/gcp"
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
		{name: "test", error: true, desc: "name shorter than minimum"},
		{name: "test-me?", error: true, desc: "non supported characters"},
		{name: "-test-me", error: true, desc: "starting with hypen"},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			c, err := caravan.NewConfigFromScratch(tc.name, "gcp", "")
			if err != nil {
				t.Fatalf("unable to create config: %s\n", err)
			}
			_, err = gcp.New(*c)
			if err == nil && tc.error || err != nil && !tc.error {
				t.Errorf("something wen wrong: want %t but got %s", tc.error, err)
			}
		})
	}
}
