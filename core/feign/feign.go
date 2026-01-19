package feign

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var globalFeignClientInterceptors = []FeignClientInterceptor{}

func RegisterGlobalFeignClientInterceptor(interceptor FeignClientInterceptor) {
	globalFeignClientInterceptors = append(globalFeignClientInterceptors, interceptor)
}

func NewClient() *FeignClient {
	return &FeignClient{
		HttpClient:   &http.Client{},
		Interceptors: globalFeignClientInterceptors,
		cloned:       false,
	}
}

type FeignClientInterceptor interface {
	Do(chain *FeignClientInterceptorChain, method string, url string, urlValues url.Values, bodyValue any, options ...*RequestOptions) *FeignClient
}

type FeignClientInterceptorChain struct {
	*FeignClient
	p int
}

func (f *FeignClientInterceptorChain) Next(method string, url string, urlValues url.Values, bodyValue any, options ...*RequestOptions) *FeignClient {
	if f.p < len(f.Interceptors) {
		f.p++
		return f.Interceptors[f.p-1].Do(f, method, url, urlValues, bodyValue, options...)
	} else {
		return f.do(method, url, urlValues, bodyValue, options...)
	}
}

func (f *FeignClientInterceptorChain) Exit(err error) *FeignClient {
	f.Error = err
	return f.FeignClient
}

/**
 *	Skip剩余的操作
 */
func (f *FeignClientInterceptorChain) Done() *FeignClient {
	return f.FeignClient
}

type FeignClient struct {
	HttpClient   *http.Client
	Interceptors []FeignClientInterceptor
	Request      *http.Request
	Response     *http.Response
	Error        error
	cloned       bool
}

type RequestOptions struct {
	Headers     map[string]string
	ContentType string
}

var DefaultRequestOptions = RequestOptions{
	ContentType: "application/json;charset=UTF-8",
}

func (m *FeignClient) clone() *FeignClient {
	return &FeignClient{
		HttpClient: &http.Client{},
		cloned:     true,
	}
}

func (m *FeignClient) Fetch(data any) error {
	if !m.cloned {
		return errors.New("invalid client")
	}

	if m.Error != nil {
		return m.Error
	}

	if m.Response == nil {
		return errors.New("no response")
	}

	if resBytes, err := io.ReadAll(m.Response.Body); err != nil {
		return err
	} else {
		if m.Response.StatusCode < 200 || m.Response.StatusCode > 299 {
			return errors.New(string(resBytes))
		}

		return json.Unmarshal(resBytes, data)
	}
}

func (m *FeignClient) Raw() ([]byte, error) {
	if !m.cloned {
		return nil, errors.New("invalid client")
	}

	if m.Error != nil {
		return nil, m.Error
	}

	if m.Response == nil {
		return nil, errors.New("no response")
	}

	if resBytes, err := io.ReadAll(m.Response.Body); err != nil {
		return nil, err
	} else {
		return resBytes, nil
	}
}

func (m *FeignClient) Get(url string, urlValues url.Values, options ...*RequestOptions) *FeignClient {
	return m.Do("GET", url, urlValues, nil, options...)
}

func (m *FeignClient) Post(url string, urlValues url.Values, bodyValue any, options ...*RequestOptions) *FeignClient {
	return m.Do("POST", url, nil, bodyValue, options...)
}

func (m *FeignClient) Put(url string, urlValues url.Values, bodyValue any, options ...*RequestOptions) *FeignClient {
	return m.Do("PUT", url, urlValues, bodyValue, options...)
}

func (m *FeignClient) Patch(url string, urlValues url.Values, bodyValue any, options ...*RequestOptions) *FeignClient {
	return m.Do("PATCH", url, urlValues, bodyValue, options...)
}

func (m *FeignClient) Delete(url string, urlValues url.Values, options ...*RequestOptions) *FeignClient {
	return m.Do("DELETE", url, urlValues, nil, options...)
}

func (m *FeignClient) Do(method string, url string, urlValues url.Values, bodyValue any, options ...*RequestOptions) *FeignClient {
	if !m.cloned {
		return m.clone().Do(method, url, urlValues, bodyValue, options...)
	}

	return (&FeignClientInterceptorChain{
		FeignClient: m,
	}).Next(method, url, urlValues, bodyValue, options...)
}

func (m *FeignClient) do(method string, url string, urlValues url.Values, bodyValue any, options ...*RequestOptions) *FeignClient {
	if m.Error != nil {
		return m
	}

	var bodyValueBytes []byte = make([]byte, 0)
	if bodyValue != nil {
		bytes, err := json.Marshal(bodyValue)
		if err != nil {
			m.Error = err
			return m
		}
		bodyValueBytes = bytes
	}

	req, err := http.NewRequest(method, m.formatUrl(url, urlValues), bytes.NewReader(bodyValueBytes))
	if err != nil {
		m.Error = err
		return m
	}

	m.Request = req
	return m.doRequest(req, options...)
}

func (m *FeignClient) formatUrl(url string, urlValues url.Values) string {
	sb := new(strings.Builder)
	sb.WriteString(url)
	if len(urlValues) > 0 {
		sb.WriteString("?")
		sb.WriteString(urlValues.Encode())
	}

	return sb.String()
}

func (m *FeignClient) doRequest(req *http.Request, options ...*RequestOptions) *FeignClient {
	var option = &DefaultRequestOptions
	if len(options) > 0 {
		option = options[0]
	}

	req.Header.Set("Content-Type", option.ContentType)
	for k, v := range option.Headers {
		req.Header.Set(k, v)
	}

	if res, err := m.HttpClient.Do(req); err != nil {
		m.Error = err
	} else {
		m.Response = res
	}
	return m
}
