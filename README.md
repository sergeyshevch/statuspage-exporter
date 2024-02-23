# Statuspage Exporter

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fsergeyshevch%2Fstatuspage-exporter.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fsergeyshevch%2Fstatuspage-exporter?ref=badge_shield)
![Build Status](https://github.com/sergeyshevch/statuspage-exporter/workflows/CI/badge.svg)
[![License](https://img.shields.io/github/license/sergeyshevch/statuspage-exporter)](/LICENSE)
[![Release](https://img.shields.io/github/release/sergeyshevch/statuspage-exporter.svg)](https://github.com/sergeyshevch/statuspage-exporter/releases/latest)
[![Docker](https://img.shields.io/docker/pulls/sergeykons/statuspage-exporter)](https://hub.docker.com/r/sergeykons/statuspage-exporter)

Statuspage exporter exports metrics from given statuspages as prometheus metrics.

Statuspage exporter is a multi-target exporter. It scrape statuspages using probes. You can read more about it
here: [Prometheus Docs](https://prometheus.io/docs/guides/multi-target-exporter/#understanding-and-using-the-multi-target-exporter-pattern)

## Supported statuspage engines:

- Statuspage.io (Widely used statuspage engine. For example by [GitHub](https://www.githubstatus.com)). You can check
  that statuspage is supported by this engine by checking that it has
  a [/api/v2/components.json](https://www.githubstatus.com/api/v2/components.json) endpoint.
- Status.io (Widely used statuspage engine. For example by [Gitlab.com](https://status.gitlab.com). You can check that
  statuspage is supported by this engine by checking footer of the page. It should contain status.io text)

Statuspage exporter will automatically detect, which engine used by statuspage and will use appropriate parser.
If this statuspage is not supported by any of the engines, then statuspage exporter will return an error.

## Some popular statuspages

### Statuspage.io based

- [GitHub](https://www.githubstatus.com)
- [Atlassian (Jira/Confluence/etc)](https://status.atlassian.com/)

### Status.io based

- [Gitlab.com](https://status.gitlab.com/)
- [Docker](https://status.docker.com/)
- [Twitter](https://status.twitterstat.us/)

## Status mapping

Different statuspage engines have different statuses. Statuspage exporter will map statuses to the same values for all
statuspages.

Special cases:
- Overall status for Status.io engine can be wrong because of statuspage text lacks information about some statuses
- Status.io statuses is parsed correctly only if it wasn't customized in status.io dashboard. Status.io doesn't provide public API so exporter relies on page text.

| Statuspage.io        | Status.io                  | Statuspage Exporter | Description                           |
|----------------------|----------------------------|---------------------|---------------------------------------|
| -                    | -                          | 0                   | Status unknown                        |
| operational          | Operational                | 1                   | System / Component operation normally |
| -                    | Planned Maintenance        | 2                   | Planned maintenance                   |
| degraded_performance | Degraded Performance       | 3                   | Degraded performance                  |
| partial_outage       | Partial Service Disruption | 4                   | Partial outage                        |
| major_outage         | Service Disruption         | 5                   | Major outage                          |
| -                    | Security Issue             | 6                   | Security incident                     |

## Running exporter

You can run the exporter with docker, kubernetes, or just as a binary. After running you can get results by http:

```bash
curl http://localhost:9747/probe?target=https://www.githubstatus.com
```

### Docker

Docker images available in Github Registry/DockerHub in all arch (amd64, arm64, arm/v7) for linux. Please be careful
with DockerHun pull limits.

| Registry        | Repository                                                                                                                         |
|-----------------|------------------------------------------------------------------------------------------------------------------------------------|
| Github Registry | [ghcr.io/sergeyshevch/statuspage-exporter](https://github.com/sergeyshevch/statuspage-exporter/pkgs/container/statuspage-exporter) |
| DockerHub       | [sergeykons/statuspage-exporter](https://hub.docker.com/r/sergeykons/statuspage-exporter)                                          |

```bash
docker run -p 9747:8080 ghcr.io/sergeyshevch/statuspage-exporter --statuspages=https://www.githubstatus.com, https://https://jira-software.status.atlassian.com
```

### Helm

```shell
helm add sergeyshevch sergeyshevch.github.io/charts
helm install sergeyshevch/statuspage-exporter --namespace statuspage-exporter --create-namespace --set statuspages[0]=https://www.githubstatus.com
```

### Binary

Please select latest available release
from [releases page](https://github.com/sergeyshevch/statuspage-exporter/releases)

```
wget https://github.com/sergeyshevch/statuspage-exporter/releases/download/v1.2.0/statuspage-exporter_v1.2.0_darwin_amd64 -O statuspage-exporter
sudo chmod +x statuspage-exporter
./statuspage-exporter
```

## Configuration

You can provide configuration using configuration file or environment variables.

Configuration file must be named as '.statuspage-exporter.yaml' and should be placed in the home directory or same
directory as the binary.

Environment variable names are the same as configuration file keys but in upper case and with underscores instead of
dots.

You can read defaults from [config.go](/pkg/config/config.go)

### Configuration file example

```yaml
http_port: 9747
# Timeout for the http client
client_timeout: 2
# Count of retries for the http client
retry_count: 3
```

## Prometheus configuration

Statuspage exporter implements the multi-target exporter pattern, so we advice
to read the guide [Understanding and using the multi-target exporter pattern
](https://prometheus.io/docs/guides/multi-target-exporter/) to get the general
idea about the configuration.

The statuspage exporter needs to be passed the target as a parameter, this can be
done with relabelling.

Example config:
```yml
scrape_configs:
  - job_name: 'statuspage'
    metrics_path: /probe
    static_configs:
      - targets:
        - https://www.githubstatus.com    # Target to probe with http.
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: 127.0.0.1:9747  # The statuspage exporter's real hostname:port.
```

You also can use Prometheus operator kind:Probe or VictoriaMetrics operator kind:VMProbe Custom resources for same purpose.

```yaml
apiVersion: monitoring.coreos.com/v1
kind: Probe
metadata:
  name: statuspage-probe
spec:
  module: http_2xx
  targets:
    staticConfig:
      static:
        - 'www.githubstatus.com'
```

## Metrics Example

```
# HELP service_status_fetch_duration_seconds Returns how long the service status fetch took to complete in seconds
# TYPE service_status_fetch_duration_seconds gauge
service_status_fetch_duration_seconds{status_page_url="https://www.githubstatus.com"} 1.078459208
# HELP statuspage_component Status of a service component. 0 - Unknown, 1 - Operational, 2 - Planned Maintenance, 3 - Degraded Performance, 4 - Partial Outage, 5 - Major Outage, 6 - Security Issue
# TYPE statuspage_component gauge
statuspage_component{component="API Requests",service="GitHub",status_page_url="https://www.githubstatus.com"} 0
statuspage_component{component="Actions",service="GitHub",status_page_url="https://www.githubstatus.com"} 0
statuspage_component{component="Codespaces",service="GitHub",status_page_url="https://www.githubstatus.com"} 0
statuspage_component{component="Copilot",service="GitHub",status_page_url="https://www.githubstatus.com"} 0
statuspage_component{component="Git Operations",service="GitHub",status_page_url="https://www.githubstatus.com"} 0
statuspage_component{component="Issues",service="GitHub",status_page_url="https://www.githubstatus.com"} 0
statuspage_component{component="Packages",service="GitHub",status_page_url="https://www.githubstatus.com"} 0
statuspage_component{component="Pages",service="GitHub",status_page_url="https://www.githubstatus.com"} 0
statuspage_component{component="Pull Requests",service="GitHub",status_page_url="https://www.githubstatus.com"} 0
statuspage_component{component="Visit www.githubstatus.com for more information",service="GitHub",status_page_url="https://www.githubstatus.com"} 0
statuspage_component{component="Webhooks",service="GitHub",status_page_url="https://www.githubstatus.com"} 0
# HELP statuspage_overall Overall status of a service0 - Unknown, 1 - Operational, 2 - Planned Maintenance, 3 - Degraded Performance, 4 - Partial Outage, 5 - Major Outage, 6 - Security Issue
# TYPE statuspage_overall gauge
statuspage_overall{service="GitHub",status_page_url="https://www.githubstatus.com"} 0
```

## License Scan

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fsergeyshevch%2Fstatuspage-exporter.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fsergeyshevch%2Fstatuspage-exporter?ref=badge_large)



