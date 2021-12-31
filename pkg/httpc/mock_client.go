package httpc

import (
	"context"
	"errors"
	"io/ioutil"
)

type MockClient struct {
	StatusCode       int
	BodyString       string
	ResponseBodyFile string
	ErrorMessage     string
}

func (c *MockClient) DoRequest(context.Context, *RequestOptions) (*Response, error) {
	if c.ErrorMessage != "" {
		return nil, errors.New(c.ErrorMessage)
	}

	var bodyBytes []byte
	if c.BodyString != "" {
		bodyBytes = []byte(c.BodyString)
	} else if c.ResponseBodyFile != "" {
		b, err := ioutil.ReadFile(c.ResponseBodyFile)
		if err != nil {
			panic(err.Error())
		}
		bodyBytes = b
	}

	res := &Response{
		StatusCode: c.StatusCode,
		BodyBytes:  bodyBytes,
	}

	return res, nil
}
