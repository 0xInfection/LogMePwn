package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"
)

func main() {
	flag.IntVar(&maxConcurrent, "threads", 10, "Number of threads to use while scanning.")
	flag.IntVar(&delay, "delay", 0, "Delay between subsequent requests for the same host to avoid overwhelming the host.")
	flag.BoolVar(&useJson, "json", false, "Use body of type JSON in HTTP requests that can contain a body.")
	// flag.BoolVar(&randomScan, "random-scan", false, "Randomly scan IP addresses (Here we go pew pew pew).")
	flag.BoolVar(&useXML, "xml", false, "Use body of type XML in HTTP requests that can contain a body.")
	flag.StringVar(&canaryToken, "token", "", "Canary token payload to use in requests; if empty, a new token will be generated.")
	flag.StringVar(&email, "email", "", "Email to use for the receiving callback notifications.")
	flag.StringVar(&webhook, "webhook", "", "Webhook to use for receiving callback notifications.")
	flag.StringVar(&userAgent, "user-agent", "", "Custom user-agent string to use; if empty, payloads will be used.")
	flag.StringVar(&urlFile, "file", "", "Specify a file containing list of hosts to scan.")
	flag.StringVar(&commonPorts, "ports", "80,443,8080", "Comma separated list of ports to scan per target.")
	flag.StringVar(&hMethods, "methods", "GET", "Comma separated list of HTTP methods to use while scanning.")
	flag.StringVar(&hHeaders, "headers", "", "Comma separated list of HTTP headers to use; if empty a default set of headers are used.")
	flag.StringVar(&hBody, "fbody", "", "Specify a format string to use as the body of the HTTP request.")
	flag.StringVar(&customServer, "custom-server", "", "Specify a custom callback server.")

	mainUsage := func() {
		fmt.Fprint(os.Stdout, lackofart, "\n")
		fmt.Fprintf(os.Stdout, "Usage:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stdout, "\nExamples:\n")
		fmt.Fprint(os.Stdout, "  ./lmp -email alerts@testing.site 1.2.3.4 1.1.1.1:8080\n")
		fmt.Fprint(os.Stdout, "  ./lmp -token xxxxxxxxxxxxxxxxxx -fbody '<padding_here>%s<padding_here>' -headers X-Custom-Header\n")
		fmt.Fprint(os.Stdout, "  ./lmp -webhook https://webhook.testing.site -file internet-ranges.lst -ports 8000,8888\n")
		fmt.Fprint(os.Stdout, "  ./lmp -email alerts@testing.site -methods GET,POST,PUT,PATCH,DELETE 1.2.3.4:8880\n\n")
	}
	flag.Usage = mainUsage
	flag.Parse()
	//fmt.Print(lackofart, "\n")

	allTargets = flag.Args()
	if len(allTargets) < 1 && len(urlFile) < 1 && !randomScan {
		flag.Usage()
		log.Println("You need to supply at least a valid target via arguments or '-file' to scan!")
		os.Exit(1)
	}

	if len(email) < 1 && len(webhook) < 1 && len(canaryToken) < 1 {
		flag.Usage()
		log.Println("You need to supply either a email or webhook to receive notifications at!")
		os.Exit(1)
	}

	fmt.Print(lackofart, "\n\n")

	for _, port := range strings.Split(commonPorts, ",") {
		allPorts = append(allPorts, strings.TrimSpace(port))
	}

	for _, method := range strings.Split(hMethods, ",") {
		allMethods = append(allMethods, strings.TrimSpace(method))
	}

	// if custom callback server is supplied
	if len(customServer) < 1 {
		if len(canaryToken) < 1 {
			canaryToken = getToken()
		}
		xload = fmt.Sprintf(canaryTokenFormat, canaryToken)
	} else {
		xload = fmt.Sprintf(genericPayFormat, customServer)
	}

	if useJson {
		var dJson struct {
			S string `json:"s"`
		}
		dJson.S = xload
		data, _ := json.Marshal(dJson)
		dummyJSON = string(data)
	} else if useXML {
		dummyXML = fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?><Request>%s</Request>`, xload)
	}

	_, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go handleInterrupt(c, &cancel)

	tnoe := time.Now()
	log.Println("Starting scan at:", tnoe.Local().String())
	go ProcessHosts()

	initDispatcher(maxConcurrent)
	dnoe := time.Now()
	fmt.Print("\n")
	if len(canaryResp.Auth) > 0 {
		manageUrl := fmt.Sprintf("https://canarytokens.org/history?token=%s&auth=%s", canaryResp.Token, canaryResp.Auth)
		log.Printf("Visit '%s' for seeing the triggers of your payloads.", manageUrl)
	}
	log.Println("Scan finished at:", dnoe.Local().String())
	log.Println("Total time taken to scan:", time.Since(tnoe).String())
	log.Println("LogMePwn is exiting.")
}
