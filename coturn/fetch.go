package coturn

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/tvi/coturn_exporter/coturn/parser"
	"github.com/tvi/coturn_exporter/coturn/types"
)

type Fetcher interface {
	Fetch(ctx context.Context) (int, map[string]types.Session, error)
}

type simpleFetcher struct {
	initURL, loginURL, psURL string
	loginForm                string
}

func NewFetcher(username, password, endpointURL string) (*simpleFetcher, error) {
	formData := url.Values{"uname": {username}, "pwd": {password}}
	ret := &simpleFetcher{loginForm: formData.Encode()}
	base, err := url.Parse(endpointURL)
	if err != nil {
		return nil, err
	}
	ret.initURL = base.String()
	base.Path = "/logon"
	ret.loginURL = base.String()

	base.Path = "/ps"
	// TODO(tvi): Make this configurable.
	ret.psURL = base.String() + "?realm=&maxsess=256"

	return ret, nil
}

func (s *simpleFetcher) Fetch(ctx context.Context) (int, map[string]types.Session, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
	req, err := http.NewRequest("GET", s.initURL, nil)
	if err != nil {
		return 0, nil, fmt.Errorf("could create request (%v): %v", s.initURL, err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("could not request (%v): %v", s.initURL, err)
	}
	if _, err := ioutil.ReadAll(resp.Body); err != nil {
		return 0, nil, fmt.Errorf("could not read response body (%v): %v", s.initURL, err)
	}

	req, err = http.NewRequest("POST", s.loginURL, strings.NewReader(s.loginForm))
	if err != nil {
		return 0, nil, fmt.Errorf("could create request (%v): %v", s.loginURL, err)
	}
	resp, err = client.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("could not request (%v): %v", s.loginURL, err)
	}
	if _, err := ioutil.ReadAll(resp.Body); err != nil {
		return 0, nil, fmt.Errorf("could not read response body (%v): %v", s.loginURL, err)
	}

	req, err = http.NewRequest("GET", s.psURL, nil)
	if err != nil {
		return 0, nil, fmt.Errorf("could create request (%v): %v", s.psURL, err)
	}
	resp, err = client.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("could not request (%v): %v", s.psURL, err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, fmt.Errorf("could not read response body (%v): %v", s.psURL, err)
	}

	return parser.Parse(b)
}
