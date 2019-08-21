package httpc

import (
	"time"
)

//Proxy setting
type Proxy struct {
	Host     string //host name
	User     string //user name for auth. default is "" which means that doesn't use.
	Password string //password for auth.
}

//Timeout timeout setting. Default is 0 which means that doesn't use and use golang default setting.
type Timeout struct {
	Connect time.Duration //connection timeout
	Read    time.Duration //read timeout
	Write   time.Duration //write timeout
	Header  time.Duration //timeout for waiting Response Header
	Idle    time.Duration //timeout for idle
}
