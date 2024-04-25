package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
)

type Client struct {
	*http.Client
	headers http.Header
	query   url.Values
	body    io.Reader
}

func New() *Client {
	return &Client{
		Client:  &http.Client{},
		headers: make(http.Header),
		query:   make(url.Values),
		body:    nil,
	}
}

func (c *Client) WithHeaders(headers map[string]string) *Client {
	for k, v := range headers {
		c.headers.Set(k, v)
	}
	return c
}

func (c *Client) WithQuery(query map[string]string) *Client {
	for k, v := range query {
		c.query.Set(k, v)
	}
	return c
}

func (c *Client) WithBody(body []byte) *Client {
	c.body = bytes.NewReader(body)
	return c
}

func (c *Client) Get(url string) (*Response, error) {
	return c.do("GET", url)
}

func (c *Client) Post(url string) (*Response, error) {
	return c.do("POST", url)
}

func (c *Client) Put(url string) (*Response, error) {
	return c.do("PUT", url)
}

func (c *Client) Patch(url string) (*Response, error) {
	return c.do("PATCH", url)
}

func (c *Client) Delete(url string) (*Response, error) {
	return c.do("DELETE", url)
}

func (c *Client) do(method string, uri string) (*Response, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}
	u.RawQuery = c.query.Encode()

	req, err := http.NewRequest(method, u.String(), c.body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header = c.headers

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return &Response{resp}, nil
}

func (c *Client) DoAll(method string, urls []string, concurrent bool) ([]*Response, error) {
	var responses []*Response

	if concurrent {
		var wg sync.WaitGroup

		wg.Add(len(urls))

		for _, url := range urls {
			go func(url string) {
				defer wg.Done()

				resp, err := c.do(method, url)
				if err != nil {
					fmt.Println(err)
					return
				}

				responses = append(responses, resp)
			}(url)
		}

		wg.Wait()
	} else {
		for _, url := range urls {
			resp, err := c.do(method, url)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}

			responses = append(responses, resp)
		}
	}

	return responses, nil
}

type Response struct {
	*http.Response
}

func (r *Response) BodyMap(result *any) (map[string]interface{}, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	contentType := r.Header.Get("Content-Type")

	switch contentType {
	case "application/json":
		err = json.Unmarshal(body, &result)
		if err != nil {
			return nil, err
		}
	case "application/x-www-form-urlencoded":
		values, err := url.ParseQuery(string(body))
		if err != nil {
			return nil, err
		}
		result = make(map[string]interface{})
		for k, v := range values {
			result[k] = v[0]
		}
	case "text/plain":
		result = make(map[string]interface{})
		result["raw"] = string(body)
	case "text/html":
		result = make(map[string]interface{})
		result["raw"] = string(body)
	case "text/xml":
		result = make(map[string]interface{})
		result["raw"] = string(body)
	default:
		result = make(map[string]interface{})
		result["raw"] = string(body)
	}

	return result, nil
}

func 
