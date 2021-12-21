package checker

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/rs/zerolog/log"
)

type ConsulChecker struct {
	GenericChecker
}

func NewConsulChecker(u, ca string, options ...func(*GenericChecker)) (cc ConsulChecker, err error) {
	client, err := TLSClient(ca)
	if err != nil {
		return cc, err
	}
	gc := NewGenericChecker(u, client)
	for _, op := range options {
		if op != nil {
			op(gc)
		}
	}
	return ConsulChecker{
		GenericChecker: *gc,
	}, nil
}

func (cc ConsulChecker) Status(ctx context.Context) bool {
	u := "/v1/status/leader"
	return cc.GenericChecker.CheckURL(ctx, u)
}

func (cc ConsulChecker) Version(ctx context.Context) (version string) {
	u := "/ui/" + cc.Datacenter + "/services"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cc.url+u, nil)
	if err != nil {
		log.Error().Msgf("error creating request: %s", err)
		return fmt.Sprintf("error: %s\n", err)
	}
	resp, err := cc.Client.Do(req)
	if err != nil {
		log.Error().Msgf("error executing request: %s", err)
		return fmt.Sprintf("error: %s\n", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Debug().Msgf("body: %s error: %s", body, err)
		log.Error().Msgf("error reading response: %s", err)
		return ""
	}

	re := regexp.MustCompile("CONSUL_VERSION: (.*) --")
	match := re.FindStringSubmatch(string(body))
	if len(match) == 2 {
		return match[1]
	}
	return "not found"
}
