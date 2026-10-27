package roku

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sidalikharef/goremote/internal/httpclient"
	"github.com/sidalikharef/goremote/internal/validate"
	"github.com/sidalikharef/goremote/tv"
)

// Client implements tv.Driver for Roku ECP.
type Client struct {
	host   string
	port   string
	client *http.Client
}

// New creates a new Roku client.
func New(host string, timeout time.Duration) (*Client, error) {
	if net.ParseIP(host) == nil {
		return nil, fmt.Errorf("invalid host")
	}
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &Client{
		host:   host,
		port:   "8060",
		client: httpclient.New(timeout),
	}, nil
}

// Name returns driver name.
func (c *Client) Name() string { return "roku" }

// baseURL builds base URL.
func (c *Client) baseURL() string {
	return fmt.Sprintf("http://%s:%s", c.host, c.port)
}

// Key sends a keypress.
func (c *Client) Key(ctx context.Context, key string) error {
	if err := validate.KeyName(key); err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL()+"/keypress/"+url.PathEscape(key), nil)
	if err != nil {
		return err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("roku error: %s", resp.Status)
	}
	return nil
}

// Apps lists installed apps.
func (c *Client) Apps(ctx context.Context) ([]tv.App, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL()+"/query/apps", nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("roku error: %s", resp.Status)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
	if err != nil {
		return nil, err
	}
	var result struct {
		XMLName xml.Name `xml:"apps"`
		Apps    []struct {
			ID   string `xml:"id,attr"`
			Name string `xml:",chardata"`
		} `xml:"app"`
	}
	if err := xml.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	apps := make([]tv.App, 0, len(result.Apps))
	for _, a := range result.Apps {
		name := strings.TrimSpace(a.Name)
		if name == "" {
			continue
		}
		apps = append(apps, tv.App{ID: a.ID, Name: name})
	}
	return apps, nil
}

// Launch launches an app.
func (c *Client) Launch(ctx context.Context, appID string) error {
	if err := validate.AppID(appID); err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL()+"/launch/"+url.PathEscape(appID), nil)
	if err != nil {
		return err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("roku error: %s", resp.Status)
	}
	return nil
}

// Capabilities returns capabilities.
func (c *Client) Capabilities() tv.Capabilities {
	return tv.Capabilities{
		Navigation: true,
		Volume:     true,
		Power:      true,
		Media:      true,
		Apps:       true,
		TextInput:  false,
		Inputs:     false,
		Pairing:    false,
	}
}
