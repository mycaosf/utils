package soap

import (
	"bytes"
	"encoding/xml"
	"io"
	"github.com/mycaosf/utils/net/httpc"
	"net/http"
	"reflect"
)

type Client interface {
	Version() int
	Build(params Params) (io.Reader, error)
	Header() http.Header
}
type Params map[string]interface{}

type client struct {
	version int
}

func (p *client) Version() int {
	return p.version
}

func (p *client) Build(params Params) (io.Reader, error) {
	var err error
	buf := &bytes.Buffer{}
	if p.version == SOAP_VERSION_1DOT1 {
		_, err = buf.WriteString(soapHeader1Dot1)
	} else {
		_, err = buf.WriteString(soapHeader1Dot2)
	}

	if err == nil {
		if err = p.buildData(params, buf); err == nil {
			_, err = buf.WriteString(soapTail)
		}
	}

	return buf, err
}

func (p *client) buildData(params Params, buf *bytes.Buffer) (err error) {
	e := xml.NewEncoder(buf)
	e.Indent("    ", "  ")

	if err = recursiveEncode(e, params); err == nil {
		e.Flush()
	}

	return
}

func recursiveEncode(e *xml.Encoder, value interface{}) (err error) {
	v := reflect.ValueOf(value)

	switch v.Kind() {
	case reflect.Map:
		for _, k := range v.MapKeys() {
			t := xml.StartElement{
				Name: xml.Name{
					Space: "",
					Local: k.String(),
				},
			}
			if err = e.EncodeToken(t); err == nil {
				if err = recursiveEncode(e, v.MapIndex(k).Interface()); err == nil {
					te := xml.EndElement{Name: t.Name}
					err = e.EncodeToken(te)
				}
			}

			if err != nil {
				break
			}
		}

	case reflect.Slice:
		for i := 0; i < v.Len() && err == nil; i++ {
			err = recursiveEncode(e, v.Index(i).Interface())
		}
	case reflect.String:
		content := xml.CharData(v.String())
		err = e.EncodeToken(content)
	}

	return
}

func (p *client) Header() http.Header {
	if p.version == SOAP_VERSION_1DOT1 {
		return header11
	} else {
		return header12
	}
}

func NewClient(version int) Client {
	return &client{version: version}
}

func init() {
	header11.Add(httpc.HTTPHeaderContentType, contentType1Dot1)
	header11.Add("SOAPAction", `"http://xmlme.com/WebServices/GetSpeech"`)
	header12.Add(httpc.HTTPHeaderContentType, contentType1Dot2)
}

const (
	SOAP_VERSION_1DOT1 = iota
	SOAP_VERSION_1DOT2
	contentType1Dot1 = "text/xml; charset=UTF-8"
	contentType1Dot2 = "application/soap+xml;charset=UTF-8"

	soapHeader1Dot1 = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
  xmlns:xsd="http://www.w3.org/2001/XMLSchema"
  xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
`
	soapHeader1Dot2 = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
  xmlns:xsd="http://www.w3.org/2001/XMLSchema"
  xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
  <soap:Body>
`
	soapTail = `
  </soap:Body> 
</soap:Envelope>`
)

var (
	header11 = make(http.Header)
	header12 = make(http.Header)
)
