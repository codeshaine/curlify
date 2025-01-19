package request

import (
	"errors"
	"net/http"
	"net/url"
)

type Get struct {
	Url    string
	Header http.Header
	Config RequestConfig
}

func (g *Get) Do() (*http.Response, error) {
	parsedUrl, err := url.Parse(g.Url)
	if err != nil {
		return nil, errors.New("invalid URL")
	}

	req := http.Request{
		Method: "GET",
		URL:    parsedUrl,
		Header: g.Header,
	}

	client := &http.Client{
		Timeout: g.Config.Timeout,
	}

	return client.Do(&req)
}

func NewGet(url string, header http.Header) Get {
	return Get{
		Url:    url,
		Header: header,
		Config: NewRequestConfig(),
	}
}
