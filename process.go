package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

func execWorker(wg *sync.WaitGroup) {
	for job := range ProcChan {
		job.RunChecks()
	}
	wg.Done()
}

func initDispatcher(workerno int) {
	wg := new(sync.WaitGroup)
	for i := 0; i < workerno; i++ {
		wg.Add(1)
		go execWorker(wg)
	}
	wg.Wait()
}

func ProcessHosts() {
	if len(urlFile) > 0 {
		file, err := os.Open(urlFile)
		if err != nil {
			log.Println(err)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			mhost := scanner.Text()
			if cidrRex.MatchString(mhost) {
				log.Printf("Found CIDR range: %s. Generating IP addresses...", mhost)
				for _, ip := range *genCIDRxIPs(mhost) {
					for _, method := range allMethods {
						thisTime := time.Now()
						fmt.Printf("\r%d/%02d/%02d %02d:%02d:%02d Total processed: %d | Current: %s",
							thisTime.Year(), thisTime.Month(), thisTime.Day(), thisTime.Hour(),
							thisTime.Minute(), thisTime.Second(), procCount, ip)
						if !strings.Contains(ip, "://") {
							ip = fmt.Sprintf("http://%s", ip)
						}
						ProcChan <- &ProcJob{
							Host:   ip,
							Method: method,
						}
						procCount++
					}
				}
			} else {
				for _, method := range allMethods {
					thisTime := time.Now()
					fmt.Printf("\r%d/%02d/%02d %02d:%02d:%02d Total processed: %d | Current: %s",
						thisTime.Year(), thisTime.Month(), thisTime.Day(), thisTime.Hour(),
						thisTime.Minute(), thisTime.Second(), procCount, mhost)
					if !strings.Contains(mhost, "://") {
						mhost = fmt.Sprintf("http://%s", mhost)
					}
					ProcChan <- &ProcJob{
						Host:   mhost,
						Method: method,
					}
					procCount++
				}
			}
		}
		if err := scanner.Err(); err != nil {
			log.Println(err.Error())
			return
		}
	} else {
		for _, mhost := range allTargets {
			if cidrRex.MatchString(mhost) {
				log.Printf("Found CIDR range: %s. Generating IP addresses...", mhost)
				for _, ip := range *genCIDRxIPs(mhost) {
					for _, method := range allMethods {
						thisTime := time.Now()
						fmt.Printf("\r%d/%02d/%02d %02d:%02d:%02d Total processed: %d | Current: %s",
							thisTime.Year(), thisTime.Month(), thisTime.Day(), thisTime.Hour(),
							thisTime.Minute(), thisTime.Second(), procCount, ip)
						if !strings.Contains(ip, "://") {
							ip = fmt.Sprintf("http://%s", ip)
						}
						ProcChan <- &ProcJob{
							Host:   ip,
							Method: method,
						}
						procCount++
					}
				}
			} else {
				for _, method := range allMethods {
					thisTime := time.Now()
					fmt.Printf("\r%d/%02d/%02d %02d:%02d:%02d Total processed: %d | Current: %s",
						thisTime.Year(), thisTime.Month(), thisTime.Day(), thisTime.Hour(),
						thisTime.Minute(), thisTime.Second(), procCount, mhost)
					if !strings.Contains(mhost, "://") {
						mhost = fmt.Sprintf("http://%s", mhost)
					}
					ProcChan <- &ProcJob{
						Host:   mhost,
						Method: method,
					}
					procCount++
				}
			}
		}
	}
	close(ProcChan)
}

func (p *ProcJob) RunChecks() {
	wg := new(sync.WaitGroup)
	// if the user hasn't already supplied a port, we generate
	// combinations of every default port and spawn a process
	// for each port
	if strings.Count(p.Host, ":") != 2 {
		log.Println("Running for default set of ports.")
		wg.Add(len(allPorts))
		for _, port := range allPorts {
			go p.ProcessHost(port, wg)
		}
	} else {
		wg.Add(1)
		go p.ProcessHost("", wg)
	}
	wg.Wait()
}

// processes a single host and a port
func (p *ProcJob) ProcessHost(port string, wg *sync.WaitGroup) error {
	var host string
	body := make([]byte, 0)
	headers := make(map[string]string)

	if len(port) < 1 {
		host = p.Host
	} else {
		host = fmt.Sprintf("%s:%s", p.Host, port)
	}

	var dynamicPayloads []string
	sanitisedDnsName := strings.ReplaceAll(strings.Split(host, "://")[1], ".", "-")
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
