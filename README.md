![Build Status](https://github.com/pokornyIt/prometheus_data_dump/workflows/Release/badge.svg)
[![GitHub](https://img.shields.io/github/license/pokornyIt/prometheus_data_dump)](/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/pokornyIt/prometheus_data_dump)](https://goreportcard.com/report/github.com/pokornyIt/prometheus_data_dump)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/pokornyit/prometheus_data_dump?label=latest)](https://github.com/pokornyIt/prometheus_data_dump/releases/latest)

# Prometheus Data Dump
Project designed to export data from the Prometheus database.
Exports are intend for further processing in other systems that do not support 
direct integration to the Prometheus system as a data source.

Each Prometheus metric is export to a separate file. 
The data is export for a defined number of days back and can be limited to selected "jobs". 
Existing exported data is overwrite by the new export. 
A special file "metrics-meta.json" is exported, which contains a description of individual metrics. 

# Program start

The program requires the entry of selected configuration parameters for its start. 
This is mainly the address of the Prometheus server from which the data will be exported.

## Configuration file

Program has config file.
```yaml
server: prometehus.server
path: ./dump
days: 2
from: "2021-02-01 10:30"
to: "2021-02-03 14:50"
step: 10
storeDirect: true
sources:
  - instance: 'localhost.+'
    includeGo: false
``` 

- **server** - FQDN or IP address of prometheus server
- **path** - Path for store export data
- **days** - Number of day to exports (1-60)
- **from** - From date and time dump data 
- **to** - To date and time dump data 
- **step** - Step for time slice in seconds (5 - 3600), default 10
- **storeDirect** - Store dump data direct to path or create inside path new subdirectory. Subdirectory name is *yyyyMMdd-HH:mm*   
- **sources** - Array for limit data to only for instance list name.
  - **instance** - Instance name for what you can export all data
  - **excludeGo** - Include metrics name starts with '*go_*'. Default mean exclude this metrics

If **from** and **to** values are defined, the **days** value is ignored.  

## Configuration line parameters
- **--config.show** - show actual configuration and exit
- **--config.file=cfg.yml** - define config file, default is cfg.yml
- **--path=./dump** - overwrite the path defined in config file
- **--server=IP** - FQDN or IP address of prometheus server
- **--from=date** - From date and time dump data, overwrite value in config file
- **--to=date** - To date and time dump data, overwrite value in config file
- **--back=days** - Number of day to export from now back, overwrite value in config file

## Example start program
Program run with all configuration from config file named "all-in.json":
```shell
./prometheus_data_dump --config.file=all-in.yml

# short version
./prometheus_data_dump -c all-in.yml
```

Program show actual configuration:
```shell
./prometheus_data_dump --config.file=all-in.yml --config.show

# short version
./prometheus_data_dump -c all-in.yml -v
```

Program run with overwrite configuration data:  
```shell
./prometheus_data_dump --config.file "all-in.yml" --path "/tmp/dump" --from "2021-02-18 10:00" --to "2021-02-19 12:00" --server=c01.server.com --log.level=debug 

# short version
./prometheus_data_dump -c all-in.yml  -p "/tmp/dump" -f "2021-02-18 10:00" -t "2021-02-19 12:00" -s c01.server.com --log.level=debug
```


# Contribute
We welcome any contributions. Please fork the project on GitHub and open Pull Requests for any proposed changes.

Please note that we will not merge any changes that encourage insecure behaviour. If in doubt please open an Issue first to discuss your proposal. 