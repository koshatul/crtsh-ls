# crtsh-ls
[![Build Status](https://travis-ci.org/koshatul/crtsh-ls.svg?branch=master)](https://travis-ci.org/koshatul/crtsh-ls)
[![GitHub license](https://img.shields.io/github/license/koshatul/crtsh-ls)](https://github.com/koshatul/crtsh-ls/blob/master/LICENSE)

crt.sh list command line tool 

### Installation

Download from the [releases](github.com/koshatul/crtsh-ls/releases/latest) page and move into your path.

### Usage

```
$ crtsh-ls --help
<domain> is "%.github.com" to show all subdomains of github.com or "github.com" to show single domain certificate

Usage:
  crtsh-ls <domain> [flags]

Flags:
  -d, --debug              Debug output.
  -f, --format string      Output formatting (go template).
                            Possible items are IssuerCaID, IssuerName, NameValue, MinCertID, MinEntryTimestamp, NotBefore, NotAfter.
                            (default "{{padlen .NameValue 20}}\t{{.NotBefore}}\t{{.NotAfter}}")
  -h, --help               help for crtsh-ls
      --only-valid         Only display still (date) valid certificates.
  -t, --timeout duration   Request timeout. (default 50s)
      --version            version for crtsh-ls
```

### Examples

#### List all certificates for a domain
```
$ crtsh-ls %.reddit.com
*.reddit.com        	2018-08-23T23:06:16	2019-07-26T00:03:40
*.reddit.com        	2018-08-23T23:06:16	2019-07-26T00:03:40
*.reddit.com        	2018-08-16T00:00:00	2020-09-02T12:00:00
...
```

### Exclude certificates that are no longer valid
```
$ crtsh-ls --only-valid %.reddit.com
*.reddit.com        	2018-08-23T23:06:16	2019-07-26T00:03:40
*.reddit.com        	2018-08-23T23:06:16	2019-07-26T00:03:40
*.reddit.com        	2018-08-16T00:00:00	2020-09-02T12:00:00
*.reddit.com        	2018-08-17T00:00:00	2020-09-02T12:00:00
*.reddit.com        	2018-08-17T00:00:00	2020-09-02T12:00:00
*.reddit.com        	2018-08-16T00:00:00	2020-09-02T12:00:00
*.reddit.com        	2018-08-14T02:26:13	2019-07-26T00:03:40
*.reddit.com        	2018-08-14T02:26:13	2019-07-26T00:03:40
*.reddit.com        	2018-08-07T15:20:11	2019-07-26T00:03:40
*.reddit.com        	2018-08-07T15:20:11	2019-07-26T00:03:40
*.reddit.com        	2018-07-25T00:07:02	2019-07-26T00:03:40
*.reddit.com        	2018-07-25T00:07:02	2019-07-26T00:03:40
*.reddit.com        	2018-07-25T00:03:40	2019-07-26T00:03:40
alb.reddit.com      	2018-06-18T00:00:00	2019-07-18T12:00:00
alb.reddit.com      	2018-06-18T00:00:00	2019-07-18T12:00:00
ads-api.reddit.com  	2017-12-12T00:00:00	2019-01-12T12:00:00
pixel.reddit.com    	2017-02-13T00:00:00	2020-02-26T12:00:00
```

### Formatted output
Using the `--format` option the output format can be specified.

Available fields are:
* IssuerCaID
* IssuerName
* NameValue - Servername on the certificate.
* MinCertID
* MinEntryTimestamp
* NotBefore - Certificate not valid *before* this date (date of issue).
* NotAfter - Certificate not valud *after* this date (date of expiry).

```
$ crtsh-ls --format=$'{{padlen .NameValue 30}}\t{{.NotAfter}}\t[{{.IssuerCaID}}]{{.IssuerName}}' %.dev1.slack.com
2015-06-25.dev1.slack.com     	2017-06-25T12:33:41	[85]C=IL, O=StartCom Ltd., OU=Secure Digital Certificate Signing, CN=StartCom Class 2 Primary Intermediate Server CA
*.dev1.slack.com              	2017-06-25T12:33:41	[85]C=IL, O=StartCom Ltd., OU=Secure Digital Certificate Signing, CN=StartCom Class 2 Primary Intermediate Server CA
*.enterprise.dev1.slack.com   	2017-06-25T12:33:41	[85]C=IL, O=StartCom Ltd., OU=Secure Digital Certificate Signing, CN=StartCom Class 2 Primary Intermediate Server CA
*.dev1.slack.com              	2016-12-23T05:15:02	[85]C=IL, O=StartCom Ltd., OU=Secure Digital Certificate Signing, CN=StartCom Class 2 Primary Intermediate Server CA
```
