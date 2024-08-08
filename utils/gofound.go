package utils

import (
	"bytes"
	"io"
	"net/http"
)

type GofoundClient struct {
	http.Client
	addr string
}

func NewGofoundClient(addr string) *GofoundClient {
	client := new(GofoundClient)
	client.addr = addr
	return client

}

func (c *GofoundClient) GET(url string) (resp *http.Response, err error) {
	return http.Get(c.Join(url))
}

func (c *GofoundClient) POST(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	return http.Post(c.Join(url), contentType, body)
}

func (c *GofoundClient) Join(url string) string {
	var b bytes.Buffer
	b.WriteString(c.addr)
	b.WriteString(url)
	return b.String()
}
