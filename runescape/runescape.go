// Go wrapper for the Runescape API.
//
// Further api documentation found at https://runescape.wiki/w/Application_programming_interface
package runescape

import (
	"net/http"
	"net/url"
	"sync"
)

const (
	baseURL   = "https://secure.runescape.com"
	userAgent = "go-runescape"
)

// Client represents a client for interacting with the Runescape API.
type Client struct {
	clientMutex sync.Mutex
	client      *http.Client

	BaseURL   *url.URL // BaseURL is the base URL of the Runescape API.
	UserAgent string   // UserAgent is the user agent string used in HTTP requests.
}

// NewClient creates a new Runescape API client.
//
// If httpClient is nil, a default HTTP client is used.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	httpClientCopy := *httpClient

	c := &Client{client: &httpClientCopy}
	c.clientMutex.Lock()
	defer c.clientMutex.Unlock()

	c.initialize()
	return c
}

// initialize initializes the Runescape API client with default values if not set.
func (c *Client) initialize() {
	if c.client == nil {
		c.client = &http.Client{}
	}

	if c.BaseURL == nil {
		c.BaseURL, _ = url.Parse(baseURL)
	}

	if c.UserAgent == "" {
		c.UserAgent = userAgent
	}
}
