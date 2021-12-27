package main

import (
	"bufio"
	"log"
	"os"
	"strings"
	"sync"
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
					queueHosts(ip)
				}
			} else {
				queueHosts(mhost)
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
					queueHosts(ip)
				}
			} else {
				queueHosts(mhost)
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
	if p.Protocol == "SSH" {
		if strings.Count(p.Host, ":") == 0 {
			wg.Add(len(allSSHPorts))
			for _, port := range allSSHPorts {
				go p.ProcessHostSSH(port, wg)
			}
		} else {
			wg.Add(1)
			go p.ProcessHostSSH("", wg)
		}
	}
	if p.Protocol == "HTTP" {
		if strings.Count(p.Host, ":") != 2 {
			wg.Add(len(allHTTPPorts))
			for _, port := range allHTTPPorts {
				go p.ProcessHostHTTP(port, wg)
			}
		} else {
			wg.Add(1)
			go p.ProcessHostHTTP("", wg)
		}
	}
	if p.Protocol == "IMAP" {
		if strings.Count(p.Host, ":") == 0 {
			wg.Add(len(allIMAPPorts))
			for _, port := range allIMAPPorts {
				go p.ProcessHostIMAP(port, wg)
			}
		} else {
			wg.Add(1)
			go p.ProcessHostIMAP("", wg)
		}
	}
	if p.Protocol == "FTP" {
		if strings.Count(p.Host, ":") == 0 {
			wg.Add(len(allFTPPorts))
			for _, port := range allFTPPorts {
				go p.ProcessHostFTP(port, wg)
			}
		} else {
			wg.Add(1)
			go p.ProcessHostFTP("", wg)
		}
	}
	wg.Wait()
}
