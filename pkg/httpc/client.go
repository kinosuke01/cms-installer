package httpc

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type Config struct {
	Scheme      string
	Host        string
	BasePath    string
	BaseHeaders map[string]string
}

type Client struct {
	Config
	client *http.Client
}

type RequestOptions struct {
	Method     string
	Path       string
	Queries    map[string]string
	Headers    map[string]string
	BodyValues url.Values
}

type Response struct {
	StatusCode int
	BodyBytes  []byte
}

func New(cnf *Config) *Client {
	return &Client{
		Config: *cnf,
		client: http.DefaultClient,
	}
}

func (c *Client) NewRequest(ctx context.Context, opts *RequestOptions) (*http.Request, error) {
	// set method
	method := opts.Method
	if method == "" {
		method = http.MethodGet
	}

	// set request url
	pReqURL, err := url.Parse(c.Scheme + "://" + c.Host)
	if err != nil {
		return nil, err
	}
	reqURL := *pReqURL
	reqURL.Path = path.Join(c.BasePath, opts.Path)
	if len(opts.Queries) > 0 {
		q := reqURL.Query()
		for k, v := range opts.Queries {
			q.Add(k, v)
		}
		reqURL.RawQuery = q.Encode()
	}

	// set request body
	bodyData := ""
	if len(opts.BodyValues) > 0 {
		bodyData = opts.BodyValues.Encode()
	}

	// build request
	req, err := http.NewRequest(method, reqURL.String(), strings.NewReader(bodyData))
	if err != nil {
		return nil, err
	}

	// set headers
	if len(c.BaseHeaders) > 0 {
		for k, v := range c.BaseHeaders {
			req.Header.Set(k, v)
		}
	}
	if len(opts.Headers) > 0 {
		for k, v := range opts.Headers {
			req.Header.Set(k, v)
		}
	}
	if bodyData != "" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	req = req.WithContext(ctx)

	return req, nil
}

func (c *Client) DoRequest(ctx context.Context, opts *RequestOptions) (*Response, error) {
	req, err := c.NewRequest(ctx, opts)
	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		StatusCode: res.StatusCode,
		BodyBytes:  bodyBytes,
	}, nil
}
