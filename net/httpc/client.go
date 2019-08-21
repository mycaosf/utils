package httpc

import (
	"encoding/base64"
	"github.com/golang/glog"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

//Get http get command, header, timeout may be il
func Get(url string, header http.Header, timeout *Timeout) (*http.Response, error) {
	return Request(http.MethodGet, url, header, nil, timeout)
}

func Put(url string, header http.Header, body io.Reader, timeout *Timeout) (*http.Response, error) {
	return Request(http.MethodPut, url, header, body, timeout)
}

func Post(url string, header http.Header, body io.Reader, timeout *Timeout) (*http.Response, error) {
	return Request(http.MethodPost, url, header, body, timeout)
}

//PutForm ContentType is not included in header. It is added auto.
func PutForm(url string, header http.Header, data url.Values, timeout *Timeout) (*http.Response, error) {
	return httpForm(http.MethodPut, url, header, data, timeout)
}

//PostForm ContentType is not included in header. It is added auto.
func PostForm(url string, header http.Header, data url.Values, timeout *Timeout) (*http.Response, error) {
	return httpForm(http.MethodPost, url, header, data, timeout)
}

func httpForm(method, url string, header http.Header, data url.Values, timeout *Timeout) (*http.Response, error) {
	if header == nil {
		header = make(http.Header)
	}
	header.Set(HTTPHeaderContentType, ContentTypeForm)

	return Request(method, url, header, strings.NewReader(data.Encode()), timeout)
}

//Request require a http command.
func Request(method, url string, header http.Header, body io.Reader, timeout *Timeout) (*http.Response, error) {
	client := &http.Client{Transport: CreateTransport(timeout, false)}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	if ProxySetting.Host != "" && ProxySetting.User != "" {
		if header == nil {
			header = make(http.Header)
		}

		auth := ProxySetting.User + ":" + ProxySetting.Password
		basic := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
		header.Set("Proxy-Authorization", basic)
	}

	req.Header = header
	if glog.V(2) {
		dump, e := httputil.DumpRequest(req, true)
		if e != nil {
			glog.Info(e)
		}
		glog.Infof("HTTP Request:\n%v", string(dump))
	}

	resp, err := client.Do(req)
	if glog.V(2) {
		dump, e := httputil.DumpResponse(resp, true)
		if e != nil {
			glog.Info(e)
		}
		glog.Infof("HTTP Response:\n%v", string(dump))
	}

	return resp, err
}
