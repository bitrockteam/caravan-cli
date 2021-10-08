package checker_test

import (
	"caravan-cli/cli/checker"
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

func TestGenericStatus(t *testing.T) {
	client := func(ch *checker.GenericChecker) {
		ch.Client = NewTestClient(func(req *http.Request) *http.Response {
			res := http.Response{
				StatusCode: 200,
			}
			return &res
		})
	}

	gc := checker.NewGenericChecker("url", client)
	res := gc.CheckURL(context.Background(), "")
	if !res {
		t.Errorf("got %t but wanted %t\n", res, true)
	}
}

func TestCheck(t *testing.T) {
	ctx := context.Background()

	type test struct {
		name       string
		status     bool
		statusCode int
		version    string
		body       string
	}

	tests := []test{
		{name: "nomad", statusCode: 200, status: true, version: "missing endpoint", body: ""},
		{name: "nomad", statusCode: 400, status: false, version: "missing endpoint", body: ""},
		{name: "consul", statusCode: 200, status: true, version: "1.2.3", body: "blah blah CONSUL_VERSION: 1.2.3 --- zzzz"},
		{name: "consul", statusCode: 500, status: false, version: "1.2.3", body: "blah blah CONSUL_VERSION: 1.2.3 --- zzzz"},
		{name: "vault", statusCode: 200, status: true, version: "1.2.4", body: `{ "Version": "1.2.4" }`},
		{name: "vault", statusCode: 500, status: false, version: "1.2.4", body: `{ "Version": "1.2.4" }`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := func(ch *checker.GenericChecker) {
				ch.Client = NewTestClient(func(req *http.Request) *http.Response {
					res := http.Response{
						StatusCode: tc.statusCode,
						Body:       io.NopCloser(strings.NewReader(tc.body)),
					}
					return &res
				})
			}
			var check checker.Checker
			switch tc.name {
			case "nomad":
				check, _ = checker.NewNomadChecker("url", "testdata/ca.empty", client)
			case "consul":
				check, _ = checker.NewConsulChecker("url", "testdata/ca.empty", client)
			case "vault":
				check, _ = checker.NewVaultChecker("url", "testdata/ca.empty", client)
			default:
				t.Fatalf("unsupported checker: %s\n", tc.name)
			}
			status := check.Status(ctx)
			version := check.Version(ctx)
			if status != tc.status {
				t.Errorf("got %t but wanted %t\n", status, true)
			}
			if version != tc.version {
				t.Errorf("got %s but wanted %s\n", version, tc.version)
			}
		})
	}
}
