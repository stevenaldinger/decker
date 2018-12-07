# Steps

1. Run `make docker_run_dvwa`
2. Visit [http://127.0.0.1:8080](http://127.0.0.1:8080) in your browser.
3. Click `Create / Reset Database`. At this point DVWA is ready to go.
4. Start an interactive terminal session with the `decker` container: `docker exec -it decker bash`
5. Run `decker /go/bin/examples/damn-vulnerable-web-app.hcl`
6. Open a new terminal tab and open another interactive terminal session with the `decker` container: `docker exec -it decker bash`
7. View plugin output reports. `ls -al /tmp/reports` will show you the reports that have finished so far.
8. One of the longest plugins that runs will be `w3af`. You can watch its output with this command: `tail -f /tmp/reports/output-w3af.txt`. A less verbose output with only the important things will be available for viewing at `/tmp/reports/w3af.output.txt` when it's finished.

This configuration can take a long time (~40 minutes) to finish because it's running a lot of scans (particularly `w3af`), but by the end you should have a lot of information to view with very little work.

# Some of the results

## nslookup

```
Server:		127.0.0.11
Address:	127.0.0.11#53

Non-authoritative answer:
Name:	dvwa
Address: 172.21.0.3
```

## nmap

```
Nmap scan report for dvwa.deployments_pen-testing (172.21.0.3)
Host is up (0.000018s latency).
Not shown: 998 closed ports
PORT     STATE SERVICE
80/tcp   open  http
3306/tcp open  mysql
MAC Address: xx:yy:zz:aa:bb:cc (Unknown)

Nmap done: 1 IP address (1 host up) scanned in 0.38 seconds
```

## nikto

```
- Nikto v2.1.6
---------------------------------------------------------------------------
+ Target IP:          172.21.0.3
+ Target Hostname:    dvwa
+ Target Port:        80
+ Start Time:         2018-12-07 14:04:23 (GMT0)
---------------------------------------------------------------------------
+ Server: Apache/2.4.7 (Ubuntu)
+ Cookie PHPSESSID created without the httponly flag
+ Retrieved x-powered-by header: PHP/5.5.9-1ubuntu4.20
+ The anti-clickjacking X-Frame-Options header is not present.
+ The X-XSS-Protection header is not defined. This header can hint to the user agent to protect against some forms of XSS
+ The X-Content-Type-Options header is not set. This could allow the user agent to render the content of the site in a different fashion to the MIME type
+ Root page / redirects to: login.php
+ Server leaks inodes via ETags, header found with file /robots.txt, fields: 0x1a 0x52156c6a290c0
+ Apache/2.4.7 appears to be outdated (current is at least Apache/2.4.12). Apache 2.0.65 (final release) and 2.2.29 are also current.
+ OSVDB-3268: /config/: Directory indexing found.
+ /config/: Configuration information may be available remotely.
+ OSVDB-3268: /docs/: Directory indexing found.
+ OSVDB-3233: /icons/README: Apache default file found.
+ /login.php: Admin login page/section found.
+ 8184 requests: 0 error(s) and 12 item(s) reported on remote host
+ End Time:           2018-12-07 14:04:35 (GMT0) (12 seconds)
---------------------------------------------------------------------------
+ 1 host(s) tested
```

## w3af

```
The server header for the remote web server is: "Apache/2.4.7 (Ubuntu)". This information was found in the request with id 38.
```

```
The x-powered-by header for the target HTTP server is "PHP/5.5.9-1ubuntu4.20". This information was found in the request with id 38.
```

```
The URL "http://dvwa/" returned an HTTP response without the recommended HTTP header X-Content-Type-Options.  This information was found in the request with id 19.
```

```
The remote Web server sent a strange HTTP response code: "405" with the message: "Method Not Allowed", manual inspection is recommended. This information was found in the request with id 42.
```

```
dir_file_brute plugin found "http://dvwa/cgi-bin/" with HTTP response code 403 and Content-Length: 279.
New URL found by dir_file_bruter plugin: "http://dvwa/cgi-bin/"
dir_file_brute plugin found "http://dvwa/docs/" with HTTP response code 200 and Content-Length: 1128.
New URL found by dir_file_bruter plugin: "http://dvwa/docs/"
dir_file_brute plugin found "http://dvwa/external/" with HTTP response code 200 and Content-Length: 1130.
New URL found by dir_file_bruter plugin: "http://dvwa/external/"
dir_file_brute plugin found "http://dvwa/config/" with HTTP response code 200 and Content-Length: 941.
dir_file_brute plugin found "http://dvwa/vulnerabilities/" with HTTP response code 200 and Content-Length: 3311.
New URL found by robots_txt plugin: "http://dvwa/robots.txt"
New URL found by dir_file_bruter plugin: "http://dvwa/config/"
New URL found by dir_file_bruter plugin: "http://dvwa/vulnerabilities/"
```

```
An HTTP response matching the web backdoor signature "cmd.jsp" was found at: "http://dvwa/cmd.jspx"; this could indicate that the server has been compromised. This vulnerability was found in the request with id 10415.
An HTTP response matching the web backdoor signature "cmd.jsp" was found at: "http://dvwa/cmd.jsp"; this could indicate that the server has been compromised. This vulnerability was found in the request with id 10443.
New URL found by find_backdoors plugin: "http://dvwa/cmd.jspx"
New URL found by find_backdoors plugin: "http://dvwa/cmd.jsp"
New URL found by web_spider plugin: "http://dvwa/login.php"
```

```
pykto plugin found a vulnerability at URL: "http://dvwa/config/". Vulnerability description: "Configuration information may be available remotely.". This vulnerability was found in the request with id 11785.
pykto plugin found a vulnerability at URL: "http://dvwa/config/". Vulnerability description: "Directory indexing found.". This vulnerability was found in the request with id 12040.
pykto plugin found a vulnerability at URL: "http://dvwa/docs/". Vulnerability description: "Directory indexing found.". This vulnerability was found in the request with id 13429.
pykto plugin found a vulnerability at URL: "http://dvwa/icons/README". Vulnerability description: "Apache default file found.". This vulnerability was found in the request with id 14083.
pykto plugin found a vulnerability at URL: "http://dvwa/login.php". Vulnerability description: "Admin login page/section found.". This vulnerability was found in the request with id 16534.
```

```
A robots.txt file was found at: "http://dvwa/robots.txt", this file might expose private URLs and requires a manual review. The scanner will add all URLs listed in this files to the analysis queue. This information was found in the request with id 62.
New URL found by robots_txt plugin: "http://dvwa/"
New URL found by pykto plugin: "http://dvwa/icons/README"
pykto plugin found a vulnerability at URL: "http://dvwa/server-status". Vulnerability description: "Apache server-status interface found (protected/forbidden)". This vulnerability was found in the request with id 17238.
New URL found by pykto plugin: "http://dvwa/server-status"
DAV seems to be correctly configured and allowing you to use the PUT method but the directory does not have the right permissions that would allow the web server to write to it. This error was found at: "http://dvwa/cgi-bin/dvXmp". This information was found in the requests with ids 17233 and 17268.
```

```
ReDoS was found at: "http://dvwa/icons/", using HTTP method GET. The modified parameter was the URL filename, with value: "". This vulnerability was found in the requests with ids 23201, 23241, 23247, 23255 and 23261.
```

```
Found 24 URLs and 24 different injections points.
The URL list is:
- http://dvwa/
- http://dvwa/cgi-bin/
- http://dvwa/cgi-bin/cmd.jsp
- http://dvwa/cgi-bin/cmd.jspx
- http://dvwa/cmd.jsp
- http://dvwa/cmd.jspx
- http://dvwa/config/
- http://dvwa/config/cmd.jsp
- http://dvwa/config/cmd.jspx
- http://dvwa/docs/
- http://dvwa/docs/cmd.jsp
- http://dvwa/docs/cmd.jspx
- http://dvwa/external/
- http://dvwa/external/cmd.jsp
- http://dvwa/external/cmd.jspx
- http://dvwa/icons/README
- http://dvwa/icons/cmd.jsp
- http://dvwa/icons/cmd.jspx
- http://dvwa/login.php
- http://dvwa/robots.txt
- http://dvwa/server-status
- http://dvwa/vulnerabilities/
- http://dvwa/vulnerabilities/cmd.jsp
- http://dvwa/vulnerabilities/cmd.jspx

```

```
The application has no protection against Click-Jacking attacks. All the received HTTP responses were found to be vulnerable, only the first 25 samples were captured as proof. This vulnerability was found in the requests with ids 19, 42, 62, 95 to 97, 379, 597, 1051, 1703, 2781, 10528 to 10530, 10534, 10538 to 10539, 10555 to 10556, 10562, 10566, 10571, 10574 to 10575, 10577, 10580, 10590 and 10592.
```
