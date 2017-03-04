package upyun

import (
	"net"
	"net/http"
	"time"
)

const (
	version = "2.1.0"

	defaultChunkSize      = 32 * 1024
	defaultConnectTimeout = time.Second * 60
)

type UpYunConfig struct {
	Bucket    string
	Operator  string
	Password  string
	Secret    string // deprecated
	Hosts     map[string]string
	UserAgent string
}

type UpYun struct {
	UpYunConfig
	httpc      *http.Client
	deprecated bool
}

func NewUpYun(config *UpYunConfig) *UpYun {
	up := &UpYun{}
	up.Bucket = config.Bucket
	up.Operator = config.Operator
	up.Password = md5Str(config.Password)
	up.Secret = config.Secret
	up.Hosts = config.Hosts
	if config.UserAgent != "" {
		up.UserAgent = config.UserAgent
	} else {
		up.UserAgent = makeUserAgent(version)
	}

	up.httpc = &http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (c net.Conn, err error) {
				return net.DialTimeout(network, addr, defaultConnectTimeout)
			},
		},
	}

	return up
}

func (up *UpYun) SetHTTPClient(httpc *http.Client) {
	up.httpc = httpc
}

func (up *UpYun) UseDeprecatedApi() {
	up.deprecated = true
}
