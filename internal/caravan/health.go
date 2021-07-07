package caravan

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

type Checker interface {
	Check() bool
	Version() string
}

type Health struct {
	url    string
	caFile string
}

func NewHealth(u, ca string) Health {
	return Health{
		url:    u,
		caFile: ca,
	}
}

func (h Health) Check() bool {
	resp, err := Get(h.url, h.caFile)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}

type VaultHealth struct {
	url    string
	caFile string
}

type VaultResponse struct {
	Version string `json:",omitempty"`
}

func NewVaultHealth(u, ca string) VaultHealth {
	return VaultHealth{
		url:    u + "v1/sys/health",
		caFile: ca,
	}
}

func (v VaultHealth) Check() string {
	resp, err := Get(v.url, v.caFile)
	if err != nil {
		return "error"
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		return "active"
	}
	if resp.StatusCode == 429 {
		return "standby"
	}
	return "none"
}

func (v VaultHealth) Version() string {
	resp, err := Get(v.url, v.caFile)
	if err != nil {
		fmt.Printf("error: %s", err)
		return ""
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("body: %s error: %s", body, err)
		return ""
	}

	var r VaultResponse
	err = json.Unmarshal(body, &r)
	if err != nil {
		fmt.Printf("body: %s error: %s\n", body, err)
		return ""
	}
	return r.Version
}

type ConsulHealth struct {
	url    string
	caFile string
}

func NewConsulHealth(u, ca string) ConsulHealth {
	return ConsulHealth{
		// TODO use better endpoint when available
		url:    u + "ui/aws-dc/services",
		caFile: ca,
	}
}

func (c ConsulHealth) Check() bool {
	resp, err := Get(c.url, c.caFile)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}

func (c ConsulHealth) Version() string {
	// TODO make more robust (change endpoint)
	resp, err := Get(c.url, c.caFile)
	if err != nil {
		fmt.Printf("error: %s", err)
		return ""
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error: %s", err)
		return ""
	}

	re := regexp.MustCompile("CONSUL_VERSION: (.*) --")
	match := re.FindStringSubmatch(string(body))
	if len(match) == 2 {
		return match[1]
	}
	return "not found"
}

type NomadHealth struct {
	url    string
	caFile string
}

func NewNomadHealth(u, ca string) NomadHealth {
	return NomadHealth{
		// TODO use better endpoint when available
		url:    u + "v1/sys/leader",
		caFile: ca,
	}
}

func (n NomadHealth) Check() bool {
	resp, err := Get(n.url, n.caFile)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}

func (n NomadHealth) Version() string {
	// TODO find endpoint
	return "not found"
}

func Get(url, ca string) (resp *http.Response, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return resp, err
	}

	// Load CA cert
	caCert, err := ioutil.ReadFile(ca)
	if err != nil {
		return resp, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		RootCAs:    caCertPool,
	}
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}

	client := &http.Client{Transport: transport}
	return client.Do(req)
}
