package httpc

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"github.com/golang/glog"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

//Get http get command, header, timeout may be il
func Get(url string, header http.Header, timeout *Timeout) (*http.Response, error) {
	return Request(http.MethodGet, url, header, nil, timeout)
}

func GetBytes(url string, header http.Header, timeout *Timeout) ([]byte, error) {
	return RequestBytes(http.MethodGet, url, header, nil, timeout)
}

func GetJson(url string, header http.Header, timeout *Timeout, v interface{}) error {
	return RequestJson(http.MethodGet, url, header, nil, timeout, v)
}

func GetXml(url string, header http.Header, timeout *Timeout, v interface{}) error {
	return RequestXml(http.MethodGet, url, header, nil, timeout, v)
}

func Put(url string, header http.Header, body io.Reader, timeout *Timeout) (*http.Response, error) {
	return Request(http.MethodPut, url, header, body, timeout)
}

func PutBytes(url string, header http.Header, body io.Reader, timeout *Timeout) ([]byte, error) {
	return RequestBytes(http.MethodPut, url, header, body, timeout)
}

func PutJson(url string, header http.Header, body io.Reader, timeout *Timeout, v interface{}) error {
	return RequestJson(http.MethodPut, url, header, body, timeout, v)
}

func PutXml(url string, header http.Header, body io.Reader, timeout *Timeout, v interface{}) error {
	return RequestXml(http.MethodPut, url, header, body, timeout, v)
}

func Post(url string, header http.Header, body io.Reader, timeout *Timeout) (*http.Response, error) {
	return Request(http.MethodPost, url, header, body, timeout)
}

func PostBytes(url string, header http.Header, body io.Reader, timeout *Timeout) ([]byte, error) {
	return RequestBytes(http.MethodPost, url, header, body, timeout)
}

func PostJson(url string, header http.Header, body io.Reader, timeout *Timeout, v interface{}) error {
	return RequestJson(http.MethodPost, url, header, body, timeout, v)
}

func PostXml(url string, header http.Header, body io.Reader, timeout *Timeout, v interface{}) error {
	return RequestXml(http.MethodPost, url, header, body, timeout, v)
}

//PutForm ContentType is not included in header. It is added auto.
func PutForm(url string, header http.Header, data url.Values, timeout *Timeout) (*http.Response, error) {
	return httpForm(http.MethodPut, url, header, data, timeout)
}

func PutFormBytes(url string, header http.Header, data url.Values, timeout *Timeout) ([]byte, error) {
	return httpFormBytes(http.MethodPut, url, header, data, timeout)
}

func PutFormJson(url string, header http.Header, data url.Values, timeout *Timeout, v interface{}) error {
	return httpFormJson(http.MethodPut, url, header, data, timeout, v)
}

func PutFormXml(url string, header http.Header, data url.Values, timeout *Timeout, v interface{}) error {
	return httpFormXml(http.MethodPut, url, header, data, timeout, v)
}

//PostForm ContentType is not included in header. It is added auto.
func PostForm(url string, header http.Header, data url.Values, timeout *Timeout) (*http.Response, error) {
	return httpForm(http.MethodPost, url, header, data, timeout)
}

func PostFormBytes(url string, header http.Header, data url.Values, timeout *Timeout) ([]byte, error) {
	return httpFormBytes(http.MethodPost, url, header, data, timeout)
}

func PostFormJson(url string, header http.Header, data url.Values, timeout *Timeout, v interface{}) error {
	return httpFormJson(http.MethodPost, url, header, data, timeout, v)
}

func PostFormXml(url string, header http.Header, data url.Values, timeout *Timeout, v interface{}) error {
	return httpFormXml(http.MethodPost, url, header, data, timeout, v)
}

func httpForm(method, url string, header http.Header, data url.Values, timeout *Timeout) (*http.Response, error) {
	if header == nil {
		header = make(http.Header)
	}
	header.Set(HTTPHeaderContentType, ContentTypeForm)

	return Request(method, url, header, strings.NewReader(data.Encode()), timeout)
}

func httpFormBytes(method, url string, header http.Header, data url.Values, timeout *Timeout) ([]byte, error) {
	if header == nil {
		header = make(http.Header)
	}
	header.Set(HTTPHeaderContentType, ContentTypeForm)

	return RequestBytes(method, url, header, strings.NewReader(data.Encode()), timeout)
}

func httpFormJson(method, url string, header http.Header, data url.Values, timeout *Timeout, v interface{}) error {
	if header == nil {
		header = make(http.Header)
	}
	header.Set(HTTPHeaderContentType, ContentTypeForm)

	return RequestJson(method, url, header, strings.NewReader(data.Encode()), timeout, v)
}

func httpFormXml(method, url string, header http.Header, data url.Values, timeout *Timeout, v interface{}) error {
	if header == nil {
		header = make(http.Header)
	}
	header.Set(HTTPHeaderContentType, ContentTypeForm)

	return RequestXml(method, url, header, strings.NewReader(data.Encode()), timeout, v)
}

func Delete(url string, header http.Header, body io.Reader, timeout *Timeout) (*http.Response, error) {
	return Request(http.MethodDelete, url, header, body, timeout)
}

func DeleteBytes(url string, header http.Header, body io.Reader, timeout *Timeout) ([]byte, error) {
	return RequestBytes(http.MethodDelete, url, header, body, timeout)
}

func DeleteJson(url string, header http.Header, body io.Reader, timeout *Timeout, v interface{}) error {
	return RequestJson(http.MethodDelete, url, header, body, timeout, v)
}

func DeleteXml(url string, header http.Header, body io.Reader, timeout *Timeout, v interface{}) error {
	return RequestXml(http.MethodDelete, url, header, body, timeout, v)
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

func RequestBytes(method, url string, header http.Header, body io.Reader, timeout *Timeout) ([]byte, error) {
	if resp, err := Request(method, url, header, body, timeout); err == nil {
		return httpParseBytes(resp)
	} else {
		return nil, err
	}
}

func RequestJson(method, url string, header http.Header, body io.Reader, timeout *Timeout, v interface{}) error {
	if resp, err := Request(method, url, header, body, timeout); err == nil {
		return httpParseJson(resp, v)
	} else {
		return err
	}
}

func RequestXml(method, url string, header http.Header, body io.Reader, timeout *Timeout, v interface{}) error {
	if resp, err := Request(method, url, header, body, timeout); err == nil {
		return httpParseXml(resp, v)
	} else {
		return err
	}
}

func httpParseJson(resp *http.Response, v interface{}) (err error) {
	var data []byte

	if v == nil {
		resp.Body.Close()
	} else if data, err = httpParseBytes(resp); err == nil {
		err = json.Unmarshal(data, v)
	}

	return
}

func httpParseXml(resp *http.Response, v interface{}) (err error) {
	var data []byte

	if v == nil {
		resp.Body.Close()
	} else if data, err = httpParseBytes(resp); err == nil {
		err = xml.Unmarshal(data, v)
	}

	return
}

func httpParseBytes(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
