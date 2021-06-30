package caravan_test

import (
	"caravan/internal/caravan"
	"testing"
)

func TestConfigFromScratch(t *testing.T) {
	type tc struct {
		name       string
		provider   string
		region     string
		domain     string
		wantDomain string
	}

	tests := []tc{
		{name: "name1", provider: "aws", region: "eu-south-1", domain: "test.org", wantDomain: "test.org"},
		{name: "name2", provider: "aws", region: "", domain: "", wantDomain: "reactive-labs.io"},
	}

	for _, tc := range tests {
		c, err := caravan.NewConfigFromScratch(tc.name, tc.provider, tc.region)
		if err != nil {
			t.Fatalf("unable to create config: %s\n", err)
		}
		if c.Region != tc.region {
			t.Errorf("region not set: got %s want %s", c.Region, "eu-south-1")
		}
		if tc.domain != "" {
			err = c.SetDomain(tc.domain)
			if err != nil {
				t.Errorf("unable to set domain %s: %s", tc.domain, err)
			}
			if c.Domain != tc.wantDomain {
				t.Errorf("domain mismatch: got %s want %s", c.Domain, tc.domain)
			}
		}
		if tc.domain == "" {
			if c.Domain != "reactive-labs.io" {
				t.Errorf("domain mismatch: got %s want %s", c.Domain, tc.domain)
			}
		}
	}
}

func TestConfigFromFile(t *testing.T) {
	type tc struct {
		name       string
		provider   string
		region     string
		domain     string
		wantDomain string
	}

	tests := []tc{
		{name: "name1", provider: "aws", region: "eu-south-1", domain: "test.org", wantDomain: "test.org"},
		{name: "name2", provider: "aws", region: "", domain: "", wantDomain: "reactive-labs.io"},
	}
	for _, tc := range tests {
		c, err := caravan.NewConfigFromScratch(tc.name, tc.provider, tc.region)
		if err != nil {
			t.Fatalf("unable to create config: %s\n", err)
		}
		if tc.domain != "" {
			err := c.SetDomain(tc.domain)
			if err != nil {
				t.Errorf("unable to set domain %s: %s", tc.domain, err)
			}
		}
		if c.Domain != tc.wantDomain {
			t.Errorf("domain mismatch: want %s, got %s", tc.wantDomain, c.Domain)
		}

		err = c.SaveConfig()
		if err != nil {
			t.Errorf("unable to save config to file: %s", err)
		}
		got, err := caravan.NewConfigFromFile()
		if err != nil {
			t.Fatalf("unable to load config from file %s: %s\n", tc.name, err)
		}
		if got.Region != tc.region {
			t.Errorf("error reloading confg from file: got %s want %s", got.Name, tc.name)
		}
	}
}
