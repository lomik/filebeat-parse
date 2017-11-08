Extendable log parser for [filebeat](https://github.com/elastic/beats).

Supported log formats:

* nginx_error_log
* yum_log (/var/log/yum.log)
* [ltsv](http://ltsv.org/)

It is very easy to add new formats. See yum_log parser and tests for example.

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