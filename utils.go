package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/valyala/fasthttp"
)

func handleInterrupt(c chan os.Signal, cancel *context.CancelFunc) {
	<-c
	(*cancel)()
	log.Fatalln("Scan cancelled by user.")
}

func cookHTTPRequest(httpMethod, requri string, headers map[string]string, body []byte) *fasthttp.Request {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(requri)
	req.Header.SetConnectionClose()
	req.Header.SetMethod(httpMethod)
	if len(body) > 0 {
		req.SetBody(body)
	}
	for key, val := range headers {
		req.Header.Set(key, val)
	}
	// set a custom user agent if supplied
	if len(userAgent) > 0 {
		req.Header.SetUserAgent(userAgent)
	}
	return req
}

func getToken() string {
	log.Println("Trying to generate a new Canary Token...")
	body := fmt.Sprintf(`------WebKitFormBoundaryTTwFOEyKMZZffBne
Content-Disposition: form-data; name="type"

log4shell
------WebKitFormBoundaryTTwFOEyKMZZffBne
Content-Disposition: form-data; name="email"

%s
------WebKitFormBoundaryTTwFOEyKMZZffBne
Content-Disposition: form-data; name="webhook"

%s
------WebKitFormBoundaryTTwFOEyKMZZffBne
Content-Disposition: form-data; name="fmt"


------WebKitFormBoundaryTTwFOEyKMZZffBne
Content-Disposition: form-data; name="memo"

[LogMePwn] Log4Shell Triggered!
------WebKitFormBoundaryTTwFOEyKMZZffBne
Content-Disposition: form-data; name="clonedsite"


------WebKitFormBoundaryTTwFOEyKMZZffBne
Content-Disposition: form-data; name="sql_server_table_name"

TABLE1
------WebKitFormBoundaryTTwFOEyKMZZffBne
Content-Disposition: form-data; name="sql_server_view_name"

VIEW1
------WebKitFormBoundaryTTwFOEyKMZZffBne
Content-Disposition: form-data; name="sql_server_function_name"

FUNCTION1
------WebKitFormBoundaryTTwFOEyKMZZffBne
Content-Disposition: form-data; name="sql_server_trigger_name"

TRIGGER1
------WebKitFormBoundaryTTwFOEyKMZZffBne
Content-Disposition: form-data; name="redirect_url"


------WebKitFormBoundaryTTwFOEyKMZZffBne--`, email, webhook)
	headers := map[string]string{
		"Accept":           "application/json, text/javascript, */*; q=0.01",
		"X-Requested-With": "XMLHttpRequest",
		"Origin":           "https://canarytokens.org",
		"Accept-Encoding":  "gzip, deflate, br",
		"Accept-Language":  "en-GB,en-US;q=0.9,en;q=0.8",
		"Referer":          "https://canarytokens.org/generate",
		"Sec-Fetch-Site":   "same-origin",
		"Sec-Fetch-Mode":   "cors",
		"Sec-Fetch-Dest":   "empty",
	}
	req := cookHTTPRequest("POST", "https://canarytokens.org/generate", headers, []byte(body))
	req.Header.SetContentLength(len(body))
	req.Header.SetUserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36")
	req.Header.SetContentType("multipart/form-data; boundary=----WebKitFormBoundaryTTwFOEyKMZZffBne")
	resp := fasthttp.AcquireResponse()

	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
	}()

	if err := httpClient.Do(req, resp); err != nil {
		log.Fatalln("Failed getting a Canary Token!", err.Error())
	}
	mbody := resp.Body()
	if err := json.Unmarshal(mbody, canaryResp); err != nil {
		log.Fatalln(err.Error())
	}
	// log.Println(string(mbody))
	if canaryResp.Error != nil {
		log.Fatalf("Error getting a Canary Token: %q", canaryResp.ErrorMessage)
	}
	log.Println("Successfully obtained a token:", canaryResp.Token)
	log.Println("Writing token details to file 'canarytoken-logmepwn.json'...")
	ioutil.WriteFile("canarytoken-logmepwn.json", mbody, 0644)
	return canaryResp.Token
}

func genCIDRxIPs(cidr string) *[]string {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		log.Fatalln(err.Error())
	}
	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); incrementIP(ip) {
		ips = append(ips, ip.String())
	}
	lenIPs := len(ips)
	switch {
	case lenIPs < 2:
		return &ips

	default:
		ips = ips[1 : len(ips)-1]
		return &ips
	}
}

func incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func getPayloads() *[]string {
	var payloads []string
	if _, err := os.Stat(customPayload); err != nil {
		// user has supplied a single payload
		if errors.Is(err, os.ErrNotExist) {
			payloads = append(payloads, customPayload)
		}
	} else {
		// the file exists
		file, err := os.Open(customPayload)
		if err != nil {
			log.Fatalln(err.Error())
		}
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			payloads = append(payloads, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.Fatalln(err.Error())
		}
		file.Close()
	}
	return &payloads
}

func processPayloads() error {
	// if user has supplied a custom payload (which will include their hostname)
	if len(customPayload) > 0 {
		xload = *getPayloads()
		return nil
	}

	// if user has callback server, we use the default payload scheme
	if len(customServer) > 0 {
		xload = append(xload, fmt.Sprintf(genericPayFormat, customServer))
		return nil
	}

	// if custom callback server is not supplied, we generate a canary token
	if len(canaryToken) < 1 {
		canaryToken = getToken()
		if len(canaryResp.Auth) > 0 {
			manageUrl := fmt.Sprintf("https://canarytokens.org/history?token=%s&auth=%s", canaryResp.Token, canaryResp.Auth)
			log.Printf("Visit '%s' for seeing the triggers of your payloads.", manageUrl)
		}
	}
	xload = append(xload, fmt.Sprintf(canaryTokenFormat, canaryToken))
	return nil
}
