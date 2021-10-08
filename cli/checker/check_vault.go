package checker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

type VaultChecker struct {
	GenericChecker
}

type VaultResponse struct {
	Version string `json:",omitempty"`
}

func NewVaultChecker(u, ca string, options ...func(*GenericChecker)) (vc VaultChecker, err error) {
	client, err := TLSClient(ca)
	if err != nil {
		return vc, err
	}
	gc := NewGenericChecker(u, client)
	for _, op := range options {
		if op != nil {
			op(gc)
		}
	}
	return VaultChecker{
		GenericChecker: *gc,
	}, nil
}

func (v VaultChecker) Status(ctx context.Context) bool {
	u := "/v1/sys/health"
	return v.GenericChecker.CheckURL(ctx, u)
}

func (v VaultChecker) Version(ctx context.Context) (version string) {
	u := "/v1/sys/health"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, v.url+u, nil)
	if err != nil {
		log.Error().Msgf("error creating request: %s\n", err)
		return fmt.Sprintf("error: %s\n", err)
	}
	resp, err := v.Client.Do(req)
	if err != nil {
		log.Error().Msgf("error executing request: %s\n", err)
		return fmt.Sprintf("error: %s\n", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Msgf("body: %s error: %s", body, err)
		return fmt.Sprintf("error: %s\n", err)
	}
	r := VaultResponse{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		log.Error().Msgf("unmarshal: %s error: %s\n", body, err)
		return fmt.Sprintf("error: %s\n", err)
	}

	return r.Version
}
