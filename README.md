# Prometheus Data Dump
Dump data from prometheus to JSON files.

Program dump all metrics for selected jobs into JSON files.

## Configuration

Program has config file.
```yaml
server: prometehus.server
path: ./dump
days: 2
jobs:
  - node_exporter_zqm
  - cucm_monitor
``` 

- **server** - FQDN or IP address of prometheus server
- **path** - Path for store export data
- **days** - Number of day to exports (1-60)
- **jobs** - limit data only for target jobs. If omitted or empty mean export all jobs

## Line parameters
- **--config.show** - show actual configuration and exit
- **--config.file=cfg.yml** - define config file, defaul is cfg.yml
- **--path=./dump** - overwrite path defined in  config file
- **--server=IP** - FQDN or IP address of prometheus server
