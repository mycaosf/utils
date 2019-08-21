package httpc

import (
	"net"
	"net/http"
	"net/url"
	"time"
)

type timeoutConn struct {
	conn    net.Conn
	timeout *Timeout
}

func (p *timeoutConn) Read(b []byte) (n int, err error) {
	c := p.timeout
	if c.Read != 0 {
		p.conn.SetReadDeadline(time.Now().Add(c.Read))
	}

	n, err = p.conn.Read(b)

	return
}

func (p *timeoutConn) Write(b []byte) (n int, err error) {
	c := p.timeout
	if c.Write != 0 {
		p.conn.SetWriteDeadline(time.Now().Add(c.Write))
	}

	n, err = p.conn.Write(b)

	return
}

func (p *timeoutConn) Close() error {
	return p.conn.Close()
}

func (p *timeoutConn) LocalAddr() net.Addr {
	return p.conn.LocalAddr()
}

func (p *timeoutConn) RemoteAddr() net.Addr {
	return p.conn.RemoteAddr()
}

func (p *timeoutConn) SetDeadline(t time.Time) error {
	return p.conn.SetDeadline(t)
}

func (p *timeoutConn) SetReadDeadline(t time.Time) error {
	return p.conn.SetReadDeadline(t)
}

func (p *timeoutConn) SetWriteDeadline(t time.Time) error {
	return p.conn.SetWriteDeadline(t)
}

//CreateTransport 创建Transport, timeout可以为nil，使用ProxySetting
func CreateTransport(timeout *Timeout, keepAlive bool) (ret *http.Transport) {
	if timeout == nil {
		ret = &http.Transport{
			DisableKeepAlives: !keepAlive,
		}
	} else {
		ret = &http.Transport{
			Dial: func(netw, addr string) (conn net.Conn, err error) {
				if timeout.Connect != 0 {
					conn, err = net.DialTimeout(netw, addr, timeout.Connect)
				} else {
					conn, err = net.Dial(netw, addr)
				}

				if err == nil && (timeout.Read != 0 || timeout.Write != 0) {
					conn = &timeoutConn{conn: conn, timeout: timeout}
				}

				return
			},

			ResponseHeaderTimeout: timeout.Header,
			DisableKeepAlives:     !keepAlive,
			IdleConnTimeout:       timeout.Idle,
		}
	}

	if ProxySetting.Host != "" {
		proxyURL, err := url.Parse(ProxySetting.Host)

		if err != nil {
			ret.Proxy = http.ProxyURL(proxyURL)
		}
	}

	return ret
}

var (
	ProxySetting = Proxy{}
)
