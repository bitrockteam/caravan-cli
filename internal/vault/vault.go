package vault

import (
	"fmt"

	vault "github.com/hashicorp/vault/api"
)

type Vault struct {
	URL    string
	Token  string
	Client *vault.Logical
}

func New(url string, token string, ca string) (v Vault, err error) {
	conf := vault.DefaultConfig()

	conf.ConfigureTLS(
		&vault.TLSConfig{
			CACert: ca,
		})
	conf.Address = url

	c, err := vault.NewClient(conf)
	if err != nil {
		return v, fmt.Errorf("error getting vault client: %w", err)
	}
	c.SetToken(token)
	return Vault{
		URL:    url,
		Token:  token,
		Client: c.Logical(),
	}, nil
}
func (v Vault) GetToken(path string) (token string, err error) {
	s, err := v.Client.Read(path)
	if err != nil {
		return "", err
	}
	t, err := s.TokenID()
	if err != nil {
		return "", err
	}
	// fmt.Printf("secret: %v\n", t)
	return t, nil
}
