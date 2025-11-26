package httpclient

import (
	"net/http"
	"net/url"
	"time"
)

// Client HTTP客户端
type Client struct {
	httpClient *http.Client
	baseURL    string
	headers    map[string]string
	timeout    time.Duration
}

// Option 配置选项
type Option func(*Client)

// WithBaseURL 设置基础URL
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// WithHeaders 设置请求头
func WithHeaders(headers map[string]string) Option {
	return func(c *Client) {
		c.headers = headers
	}
}

// WithTimeout 设置超时时间
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.timeout = timeout
	}
}

// WithProxy 设置代理
func WithProxy(proxyURL string) Option {
	return func(c *Client) {
		if proxyURL != "" {
			proxy, err := url.Parse(proxyURL)
			if err == nil {
				transport := &http.Transport{
					Proxy: http.ProxyURL(proxy),
				}
				c.httpClient.Transport = transport
			}
		}
	}
}

// NewHTTPClient 创建HTTP客户端
func NewHTTPClient(opts ...Option) *Client {
	client := &Client{
		httpClient: &http.Client{},
		headers:    make(map[string]string),
		timeout:    30 * time.Second,
	}

	for _, opt := range opts {
		opt(client)
	}

	client.httpClient.Timeout = client.timeout

	return client
}

// Do 执行HTTP请求
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	// 设置基础URL
	if c.baseURL != "" && req.URL.Host == "" {
		baseURL, _ := url.Parse(c.baseURL)
		req.URL.Scheme = baseURL.Scheme
		req.URL.Host = baseURL.Host
	}

	// 设置默认请求头
	for k, v := range c.headers {
		if req.Header.Get(k) == "" {
			req.Header.Set(k, v)
		}
	}

	return c.httpClient.Do(req)
}

// Get GET请求
func (c *Client) Get(urlStr string) (*http.Response, error) {
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}
