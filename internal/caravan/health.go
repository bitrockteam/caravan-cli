package caravan

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"
)

type Checker interface {
	Check() bool
	Version() string
}

type VaultHealth struct {
	url string
}

type VaultResponse struct {
	Version string `json:",omitempty"`
}

func NewVaultHealth(u string) VaultHealth {
	return VaultHealth{
		url: u + "v1/sys/health",
	}
}

func (h VaultHealth) Check() string {
	resp, err := Get(h.url)
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

func (h VaultHealth) Version() string {
	resp, err := Get(h.url)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return ""
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("body: %s error: %s\n", body, err)
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
	url string
}

func NewConsulHealth(u string) ConsulHealth {
	return ConsulHealth{
		// TODO use better endpoint when available
		url: u + "ui/aws-dc/services",
	}
}

func (c ConsulHealth) Check() bool {
	resp, err := Get(c.url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return true
}

func (c ConsulHealth) Version() string {
	// TODO make more robust (change endpoint)
	resp, err := Get(c.url)
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
	url string
}

func NewNomadHealth(u string) NomadHealth {
	return NomadHealth{
		// TODO use better endpoint when available
		url: u + "v1/sys/leader",
	}
}

func (n NomadHealth) Check() bool {
	resp, err := Get(n.url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return true
}

func (n NomadHealth) Version() string {
	// TODO find endpoint
	return "not found"
}

func Get(url string) (resp *http.Response, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return resp, err
	}
	// TODO use generated certs available
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	return client.Do(req)
}
