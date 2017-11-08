# filebeat-parse 

Extendable log parser for [filebeat](https://github.com/elastic/beats).

Supported log formats:

* nginx_error_log
* yum_log (/var/log/yum.log)
* [ltsv](http://ltsv.org/)

It is very easy to add new formats. See yum_log parser and tests for example.

Sample input and output:
```
input (line in nginx_error_log):
2012/04/15 22:01:47 [error] 3258#0: *887 upstream timed out (110: Connection timed out) while reading upstream, client: 192.168.126.1, server: *.example.com, request: "GET /wp-admin/options-general.php?page=ozh_ta&action=import_all&time=1334482870&_wpnonce=89590fa285 HTTP/1.1", upstream: "fastcgi://unix:/var/run/php-fpm/php-fpm.sock:", host: "www.example.com", referrer: "http://www.example.com/wp-login.php?redirect_to=http%3A%2F%2Fwww.example.com%2Fwp-admin%2Foptions-general.php%3Fpage%3Dozh_ta%26action%3Dimport_all%26time%3D1334482870%26_wpnonce%3D89590fa285&reauth=1"

output (sent to logstash):
{
  "@timestamp": "2017-11-08T19:11:28.553Z",
  "beat": {
    "hostname": "localhost",
    "name": "localhost",
    "version": "5.4.1"
  },
  "client": "192.168.126.1",
  "http_host": "www.example.com",
  "input_type": "log",
  "level": "error",
  "message": "upstream timed out (110: Connection timed out) while reading upstream",
  "offset": 577,
  "pid": 3258,
  "referrer": "http://www.example.com/wp-login.php?redirect_to=http%3A%2F%2Fwww.example.com%2Fwp-admin%2Foptions-general.php%3Fpage%3Dozh_ta%26action%3Dimport_all%26time%3D1334482870%26_wpnonce%3D89590fa285\u0026reauth=1",
  "request": "GET /wp-admin/options-general.php?page=ozh_ta\u0026action=import_all\u0026time=1334482870\u0026_wpnonce=89590fa285 HTTP/1.1",
  "server": "*.example.com",
  "sid": 887,
  "source": "/var/log/nginx/nginx_error_log",
  "tid": 0,
  "timestamp": "2012-04-15T22:01:47.000+0400",
  "type": "log",
  "upstream": "fastcgi://unix:/var/run/php-fpm/php-fpm.sock:"
}
```

## Usage
```yaml
# enable parse processor
processors:
- parse:
  # define field which store log format (parser)
  # events without this field are skipped by parse processor
  field: "parse"

# define inputs with "parse" extra field
filebeat.prospectors:

- input_type: log
  paths:
    - /var/log/nginx/nginx_error_log
  fields:
    parse: nginx_error_log
  fields_under_root: true
  multiline.match: after
  multiline.negate: true
  multiline.pattern: '^[0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} \['

- input_type: log
  paths:
    - /var/log/yum.log
  fields:
    parse: yum_log
  fields_under_root: true
```

## Build
Required golang
```bash
# clone "module" source
git clone https://github.com/lomik/filebeat-parse.git filebeat-parse
cd filebeat-parse

# fetch fixed version of https://github.com/elastic/beats as git submodule
make submodules

# run tests and benchmarks. optional
make test

# inject new action into filebeat and build it
make

# check
./filebeat -version
```