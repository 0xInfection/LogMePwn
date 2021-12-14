# LogMePwn
A fully automated, reliable, super-fast, mass scanning and validation toolkit for the Log4J RCE CVE-2021-44228 vulnerability. With enough amount of hardware and threads, it is capable of scanning the entire internet within a day.

![image](https://user-images.githubusercontent.com/39941993/146040886-339d1095-e861-4f1c-a009-b99732462a2b.png)

## How it works?
LogMePwn works by making use of [Canary Tokens](https://canarytokens.org), which in-turn provides email and webhook notifications to your preferred communication channel. If you have a custom callback server, you can definitely use it too!

## Installation & Usage
To use the tool, you can grab a binary from the [Releases](https://github.com/0xInfection/LogMePwn/releases) section as per your distribution and use it. If you want to build the tool, you'll need Go >= 1.13. Simple clone the repo and run `go build`.

Here's the basic usage of the tool:
```groovy
$ ./lmp --help

    +---------------------+
    |   L o g M e P w n   |
    +---------------------+  v1.0

                ~ 0xInfection
Usage:
  -custom-server string
        Specify a custom callback server.
  -delay int
        Delay between subsequent requests for the same host to avoid overwhelming the host.
  -email string
        Email to use for the receiving callback notifications.
  -fbody string
        Specify a format string to use as the body of the HTTP request.
  -file string
        Specify a file containing list of hosts to scan.
  -headers string
        Comma separated list of HTTP headers to use; if empty a default set of headers are used.
  -json
        Use body of type JSON in HTTP requests that can contain a body.
  -methods string
        Comma separated list of HTTP methods to use while scanning. (default "GET")
  -ports string
        Comma separated list of ports to scan per target. (default "80,443,8080")
  -threads int
        Number of threads to use while scanning. (default 10)
  -token string
        Canary token payload to use in requests; if empty, a new token will be generated.
  -user-agent string
        Custom user-agent string to use; if empty, payloads will be used.
  -webhook string
        Webhook to use for receiving callback notifications.
  -xml
        Use body of type XML in HTTP requests that can contain a body.

Examples:
  ./lmp -email alerts@testing.site 1.2.3.4 1.1.1.1:8080
  ./lmp -token xxxxxxxxxxxxxxxxxx -fbody '<padding_here>%s<padding_here>' -headers 'X-Custom-Header'
  ./lmp -webhook https://webhook.testing.site -file internet-ranges.lst -ports 8000,8888
  ./lmp -email alerts@testing.site -methods GET,POST,PUT,PATCH,DELETE 1.2.3.4:8880
```

#### Specifying targets
The targets can be specified in two ways, via the command line interface as arguments, or via a file. Example:
```groovy
./lmp <other args here> 1.1.1.1:8080 1.2.3.4:80 1.1.2.2:443
./lmp <other args here> -file internet-ranges.lst
```
The hosts can may contain ports, if not, the set of ports mentioned in `-ports` will be considered for scanning. The default ports list are:
- 80
- 443
- 8080

#### Specifying notification channels
The notification channels can be any of the following:
- Email (`-email`)
- Webhook (`-webhook`)
- Custom DNS callback server (`-custom-server`)

The tool makes use of Canary Tokens, you can create one from [here](https://canarytokens.org/generate), or let the tool create a token for you. If the tool creates a token, that will be written to a file named `canarytoken-logmepwn.json`, which will include the token itself and the auth (both of which you'll need to view triggers via the web interface).

If you already have a token, you can use the `-token` argument to use the token directly and not create a new one.

> __NOTE:__ If you supply either an email or a webhook, the tool will create a custom canary token. If you use a custom callback server, tokens do not come into play.

#### Sending requests
The tool offers great flexibility when sending requests. By default the tool uses GET requests. A default set of headers are used, each of which contains a payload in its value. You can specify a custom set of headers via the `-headers` argument.

You can specify the list of HTTP methods to use for scanning via the `-methods` switch. For requests that contain a body, e.g. `POST`, `PUT`, etc, you can customize content of the bodies.

By default the tool sends a payload directly via the body. The tool offers customization fo the body in the following ways:
- Specify `-json` to have the request body as type JSON.
- `-xml` for XML format.
- `-fbody` to specify a custom format string where the payload will be injected. This allows complex request creation when testing. For example, if you want to send the content as HTML, it can look like this:
  ```s
  ./lmp -fbody '<html>%s</html>' -methods 'POST,PUT' 1.2.3.4
  ```

You can specify a custom user-agent header value via the `-user-agent` switch.

#### Concurrent scanning
The tool is optimized for scanning a wide range of targets. With sufficient amount of network bandwidth and hardware, you can scan the entire IPv4 space within a day. The default number of concurrent threads to use while scanning is set at just 10 (optimised for reliability on local hardware). The value can go upto thousands (I'll leave the benchmarking task upto you). :)

Use the `-threads` switch to supply the number of threads to use with the tool.

#### Specifying delay
Since a lot of HTTP requests are involved, it might be a cumbersome job for the remote host to handle the requests. The `-delay` parameter is here to help you with those cases. You can specify a delay value in seconds -- which will be used be used in between two subsequent requests to the same port on a server.

## Demo
To demo the scanner, I make use of a vulnerable setup from [@christophetd](https://twitter.com/christophetd) using docker:
```groovy
docker run -p 8080:8080 ghcr.io/christophetd/log4shell-vulnerable-app
```
![image](https://user-images.githubusercontent.com/39941993/146034544-a0c0e60d-00db-44ae-823a-5e5834888108.png)

Then I run the tool against the setup:
```groovy
./lmp -email alerts@testing.site 127.0.0.1:8080
```
![image](https://user-images.githubusercontent.com/39941993/146034732-5600761b-008e-4119-83ce-b5b0f6686b7d.png)

Which immediately triggered a few DNS lookups visible on the token history page as well as my email:

<img src="https://user-images.githubusercontent.com/39941993/146039240-0d34e4d8-284f-4377-bde3-ea13f9f7f5eb.png" width=49% /> <img src="https://user-images.githubusercontent.com/39941993/146039600-ab2a71b1-ec92-4cef-bae4-f3f46dc2ffd6.png" width=49% />
                                                                                                                                
## Ideas & future roadmap
- [ ] Built-in capability to spin up a custom DNS callback server.
- [ ] Ability to identify all probable input fields by observing a basic HTTP response.
- [ ] Obfuscation payload generation.

## License & Version
The tool is licensed under the GNU GPLv3. LogMePwn is currently at v1.0.

## Credits
Shoutout to the team at [Thinkst Canary](https://canary.tools/) for their amazing Canary Tokens project.

> Crafted with ♡ by [Pinaki (@0xInfection)](https://twitter.com/0xinfection).
