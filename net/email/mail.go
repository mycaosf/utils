package email

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/smtp"
	"net/textproto"
	"strings"
)

type EMail struct {
	Host, User, Key string
	Port            int
	From            string
	To              []string
	Subject         string
	Msg             string
	Attachments     []string
}

func EMailSend(p *EMail) error {
	c, err := p.connect()
	if err != nil {
		return err
	}

	defer c.Close()
	if err = p.header(c); err == nil {
		var w io.WriteCloser
		if w, err = p.prepareData(c); err == nil {
			if p.Attachments == nil {
				err = p.data(w)
			} else {
				err = p.dataMP(w)
			}
			w.Close()
		}
	}

	return err
}

func (p *EMail) connect() (c *smtp.Client, err error) {
	addr := fmt.Sprintf("%v:%v", p.Host, p.Port)
	c, err = dial(addr)

	return
}

func (p *EMail) header(c *smtp.Client) (err error) {
	auth := smtp.PlainAuth("", p.User, p.Key, p.Host)
	if err = c.Auth(auth); err == nil {
		if err = c.Mail(p.From); err == nil {
			for _, v := range p.To {
				if err = c.Rcpt(v); err != nil {
					break
				}
			}
		}
	}

	return err
}

func (p *EMail) prepareData(c *smtp.Client) (w io.WriteCloser, err error) {
	if w, err = c.Data(); err == nil {
		str := "To: " + p.To[0] + "\r\n"
		for _, v := range p.To[1:] {
			str += "Cc: " + v + "\r\n"
		}
		str += "From: " + p.From + "\r\n"
		str += "Subject: " + p.Subject + "\r\n"

		if _, err = w.Write([]byte(str)); err != nil {
			w.Close()
		}
	}

	return
}

func (p *EMail) data(w io.WriteCloser) (err error) {
	str := header + "\r\n" + p.Msg

	_, err = w.Write([]byte(str))

	return
}

func (p *EMail) dataMP(w io.WriteCloser) (err error) {
	wpart := multipart.NewWriter(w)
	str := headerMP + ` boundary="` + wpart.Boundary() + `"` + "\r\nContent-Language: en-US\r\n\r\n"
	defer wpart.Close()

	if _, err = w.Write([]byte(str)); err == nil {
		var part io.Writer
		var h textproto.MIMEHeader
		if len(p.Msg) > 0 {
			h = make(textproto.MIMEHeader)
			h.Add("Content-Type", "text/plain; charset=utf-8; format=flowed")
			h.Add("Content-Transfer-Encoding", "8bit")

			part, err = wpart.CreatePart(h)
			if err == nil {
				_, err = part.Write([]byte(p.Msg))
			}
		}

		if err == nil {
			for _, fileName := range p.Attachments {
				name := fileName
				if idx := strings.LastIndexAny(fileName, "\\/"); idx >= 0 {
					name = fileName[idx+1:]
				}
				h = make(textproto.MIMEHeader)
				h.Add("Content-Type", "application/octet-stream")
				h.Add("Content-Description", name)
				h.Add("Content-Transfer-Encoding", "base64")
				h.Add("Content-Disposition", `attachment;filename="`+name+`"`)
				part, err = wpart.CreatePart(h)

				if err == nil {
					var fileData []byte
					fileData, err = ioutil.ReadFile(fileName)
					if err == nil {
						b := make([]byte, base64.StdEncoding.EncodedLen(len(fileData)))
						base64.StdEncoding.Encode(b, fileData)
						_, err = part.Write(b)
					}
				}

				if err != nil {
					break
				}
			}
		}
	}

	return
}

func dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		return nil, err
	}

	Host, _, _ := net.SplitHostPort(addr)

	return smtp.NewClient(conn, Host)
}

const (
	header   = "Content-Type: text/plain; charset=utf-8; format=flowed\r\n"
	headerMP = "MIME-Version: 1.0\r\nContent-Type: multipart/mixed;\r\n"
)
