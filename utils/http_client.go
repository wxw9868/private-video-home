package utils

import (
	"bytes"
	"io"
	"net/http"
)

type HttpClient struct {
	http.Client
	addr string
}

func NewHttpClient(addr string) *HttpClient {
	client := new(HttpClient)
	client.addr = addr
	return client

}

func (c *HttpClient) GET(url string) (resp *http.Response, err error) {
	return http.Get(c.Join(url))
}

func (c *HttpClient) POST(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	return http.Post(c.Join(url), contentType, body)
}

func (c *HttpClient) Join(url string) string {
	var b bytes.Buffer
	b.WriteString(c.addr)
	b.WriteString(url)
	return b.String()
}
