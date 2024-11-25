// Copyright (c) HashiCorp, Inc.

package cronitor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	endpoint string
	apiKey   string
	client   *http.Client
}

type NewClientOpts struct {
	Endpoint string
	ApiKey   string
	Client   *http.Client
}

func NewClient(opts NewClientOpts) *Client {
	if opts.Endpoint == "" {
		opts.Endpoint = "https://cronitor.io/api"
	}
	if opts.Client == nil {
		opts.Client = http.DefaultClient
	}

	return &Client{
		endpoint: opts.Endpoint,
		apiKey:   opts.ApiKey,
		client:   opts.Client,
	}
}

func (c *Client) Get(ctx context.Context, id string) (*Monitor, error) {
	req, err := c.request(ctx, http.MethodGet, fmt.Sprintf("/monitors/%s", id), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get monitor %s: %w", id, err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: url: %s, code %d", ErrFailedGetMonitor, req.URL.String(), resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	mon := &Monitor{}
	if err := json.Unmarshal(body, mon); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return mon, nil
}

func (c *Client) request(ctx context.Context, method, endpoint string, body any) (*http.Request, error) {
	var br io.Reader
	if body != nil {
		by, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		br = bytes.NewReader(by)
	}
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", c.endpoint, endpoint), br)
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %w", err)
	}

	req = req.WithContext(ctx)
	req.SetBasicAuth(c.apiKey, "")

	return req, nil
}
