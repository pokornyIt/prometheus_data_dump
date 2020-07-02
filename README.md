![Build Status](https://github.com/pokornyIt/prometheus_data_dump/workflows/Build/badge.svg)
[![License](https://img.shields.io/github/license/pokornyIt/prometheus_data_dump)](/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/pokornyIt/prometheus_data_dump)](https://goreportcard.com/report/github.com/pokornyIt/nut_exporter)
![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/pokornyit/prometheus_data_dump?label=latest)

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
step: 10
jobs:
  - node_exporter_zqm
  - cucm_monitor
``` 

- **server** - FQDN or IP address of prometheus server
- **path** - Path for store export data
- **days** - Number of day to exports (1-60)
- **step** - Step for time slice in seconds (5 - 3600), default 10
- **jobs** - limit data only for target jobs. If omitted or empty mean export all jobs

## Configuration line parameters
- **--config.show** - show actual configuration and exit
- **--config.file=cfg.yml** - define config file, default is cfg.yml
- **--path=./dump** - overwrite path defined in  config file
- **--server=IP** - FQDN or IP address of prometheus server

# Contribute
We welcome any contributions. Please fork the project on GitHub and open Pull Requests for any proposed changes.

Please note that we will not merge any changes that encourage insecure behaviour. If in doubt please open an Issue first to discuss your proposal. 