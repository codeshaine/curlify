package request

import (
	"io"
	"net/http"
	"net/url"
)

type Post struct {
	Url    string
	Header http.Header
	Config RequestConfig
	Body   io.ReadCloser
}

func (p *Post) Do() (*http.Response, error) {
	parsedUrl, err := url.Parse(p.Url)
	if err != nil {
		panic(err)
	}

	req := &http.Request{
		Method: "POST",
		URL:    parsedUrl,
		Header: p.Header,
		Body:   p.Body,
	}

	client := &http.Client{
		Timeout: p.Config.Timeout,
	}

	return client.Do(req)
}

func NewPost(url string, header http.Header, body io.ReadCloser) Post {
	return Post{
		Url:    url,
		Header: header,
		Config: NewRequestConfig(),
		Body:   body,
	}
}
