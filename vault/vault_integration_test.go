// +build integration

package vault_test

import (
	"caravan/internal/vault"
	"testing"
)

func TestVault(t *testing.T) {
	v, err := vault.New("https://vault.test01.reactive-labs.io", "s.ArATSZWpNNGyNRyZzlC49HiM", "../../.caravan/test01/caravan-infra-aws/ca_certs.pem")
	if err != nil {
		t.Errorf("error accessing vault: %s", err)
	}
	tok, err := v.GetToken("nomad/creds/token-manager")
	if err != nil {
		t.Errorf("error getting token: %s", err)
	}
	if tok == "" {
		t.Errorf("error: found empty token")
	}
}
