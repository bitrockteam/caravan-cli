package checker

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
)

type Checker interface {
	Status(ctx context.Context) bool
	Version(ctx context.Context) string
}

// Check bundles the url and http client.
type GenericChecker struct {
	Client     *http.Client
	url        string
	Datacenter string
}

func NewGenericChecker(u string, options ...func(*GenericChecker)) *GenericChecker {
	n := GenericChecker{
		url: u,
	}
	for _, op := range options {
		if op != nil {
			op(&n)
		}
	}
	return &n
}

func (c GenericChecker) CheckURL(ctx context.Context, u string) bool {
	log.Debug().Msgf("checking URL: %s", c.url+u)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url+u, nil)
	if err != nil {
		log.Error().Msgf("error creating request: %s", err)
		return false
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		log.Error().Msgf("error executing request: %s", err)
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}

func TLSClient(ca string) (func(*GenericChecker), error) {
	if _, err := os.ReadFile(ca); err != nil {
		return nil, err
	}
	return func(gc *GenericChecker) {
		caCert, _ := os.ReadFile(ca)
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		// Setup HTTPS client
		tlsConfig := &tls.Config{
			MinVersion: tls.VersionTLS12,
			RootCAs:    caCertPool,
		}
		transport := &http.Transport{TLSClientConfig: tlsConfig}

		gc.Client = &http.Client{Transport: transport}
	}, nil
}
