package http

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type HttpRequest struct {
	Base               string
	URL                string
	Method             string
	Headers            map[string]string
	Body               interface{}
	Username, Password string
	ContentType        string
	Response           *http.Response
	StatusCode         int
	Error              error
	bytes              []byte
	executed           bool
}

type RequestParameters map[string]string

func NewHttpRequest(base ...string) *HttpRequest {
	if len(base) > 0 {
		return &HttpRequest{Base: base[0]}
	}
	return &HttpRequest{}
}

func (r *HttpRequest) GET(url string, params ...RequestParameters) *HttpRequest {
	r.Method = http.MethodGet
	r.URL = r.urlParams(url, params...)
	return r
}

func (r *HttpRequest) POST(url string, params ...RequestParameters) *HttpRequest {
	r.Method = http.MethodPost
	r.URL = r.urlParams(url, params...)
	return r
}

func (r *HttpRequest) PUT(url string, params ...RequestParameters) *HttpRequest {
	r.Method = http.MethodPut
	r.URL = r.urlParams(url, params...)
	return r
}

func (r *HttpRequest) PATCH(url string, params ...RequestParameters) *HttpRequest {
	r.Method = http.MethodPatch
	r.URL = r.urlParams(url, params...)
	return r
}

func (r *HttpRequest) DELETE(url string, params ...RequestParameters) *HttpRequest {
	r.Method = http.MethodDelete
	r.URL = r.urlParams(url, params...)
	return r
}

func (r *HttpRequest) urlParams(url string, params ...RequestParameters) string {
	if len(params) > 0 {
		for key, value := range params[0] {
			url = strings.ReplaceAll(url, "{"+key+"}", value)
		}
	}
	return url
}

func (r *HttpRequest) BODY(body interface{}) *HttpRequest {
	r.Body = body
	return r
}

func (r *HttpRequest) HEADER(head string, value string) *HttpRequest {
	r.Headers[head] = value
	return r
}

func (r *HttpRequest) USERNAME(username string) *HttpRequest {
	r.Username = username
	return r
}

func (r *HttpRequest) PASSWORD(password string) *HttpRequest {
	r.Password = password
	return r
}

func (r *HttpRequest) fullURL() (result string) {
	// r.Base + r.URL => http://username:password@host:port/path
	result = strings.TrimSuffix(r.Base, "/") + "/" + strings.TrimPrefix(r.URL, "/")
	if r.Username != "" {
		if strings.HasPrefix(result, "http://") {
			result = strings.Replace(result, "http://", "http://"+r.Username+":"+r.Password+"@", 1)
		} else if strings.HasPrefix(result, "https://") {
			result = strings.Replace(result, "https://", "https://"+r.Username+":"+r.Password+"@", 1)
		}
	}
	return result
}

func (r *HttpRequest) Do() *HttpRequest {
	if r.executed || r.Error != nil {
		return r
	}
	r.executed = true

	var (
		body io.Reader
		req  *http.Request
		err  error
	)
	if r.Body != nil {
		switch t := r.Body.(type) {
		case string:
			body = strings.NewReader(t)
		case []byte:
			body = strings.NewReader(string(t))
		default:
			if bytes, err := json.Marshal(r.Body); err == nil {
				body = strings.NewReader(string(bytes))
			}
		}
	}
	req, err = http.NewRequest(r.Method, r.fullURL(), body)
	if err != nil {
		r.Error = err
		return r
	}

	for k, v := range r.Headers {
		req.Header.Set(k, v)
	}

	var resp *http.Response
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		r.Error = err
		return r
	}
	defer resp.Body.Close()

	r.Response = resp
	r.StatusCode = resp.StatusCode
	r.executed = true

	if bytes, err := ioutil.ReadAll(resp.Body); err == nil {
		r.bytes = bytes
	} else {
		r.Error = err
	}

	return r
}

func (r *HttpRequest) Result() (status int, body []byte, err error) {
	if !r.executed {
		return r.Do().Result()
	}

	if r.Response != nil {
		status = r.Response.StatusCode
	}

	if r.Error != nil {
		err = r.Error
	}

	body = r.bytes

	return
}

func (r *HttpRequest) ResultTo(result interface{}) (status int, err error) {
	if !r.executed {
		return r.Do().ResultTo(result)
	}

	if r.Response != nil {
		status = r.Response.StatusCode
	}

	if r.Error != nil {
		err = r.Error
	}

	if len(r.bytes) > 0 {
		switch r.ContentType {
		case "application/xml":
			err = xml.Unmarshal(r.bytes, result)
		case "application/json":
		default:
			err = json.Unmarshal(r.bytes, result)
		}
	}
	return
}

func (r *HttpRequest) String() (result string, err error) {
	if !r.executed {
		return r.Do().String()
	}

	if r.Error != nil {
		err = r.Error
	}

	if len(r.bytes) > 0 {
		result = string(r.bytes)
	}

	return
}

func (r *HttpRequest) Fetch(data interface{}) (status int, err error) {
	if !r.executed {
		return r.Do().Fetch(data)
	}

	if r.Error != nil {
		err = r.Error
	}

	if len(r.bytes) > 0 {
		switch r.ContentType {
		case "application/xml":
			err = xml.Unmarshal(r.bytes, data)
		case "application/json":
		default:
			err = json.Unmarshal(r.bytes, data)
		}
	}

	if r.Response != nil {
		status = r.Response.StatusCode
	}

	return
}
