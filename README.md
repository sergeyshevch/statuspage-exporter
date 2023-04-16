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
curl http://localhost:8080/probe?target=https://www.githubstatus.com
```

### Docker

Docker images available in Github Registry/DockerHub in all arch (amd64, arm64, arm/v7) for linux. Please be careful
with DockerHun pull limits.

| Registry        | Repository                                                                                                                         |
|-----------------|------------------------------------------------------------------------------------------------------------------------------------|
| Github Registry | [ghcr.io/sergeyshevch/statuspage-exporter](https://github.com/sergeyshevch/statuspage-exporter/pkgs/container/statuspage-exporter) |
| DockerHub       | [sergeykons/statuspage-exporter](https://hub.docker.com/r/sergeykons/statuspage-exporter)                                          |

```bash
docker run -p 8080:8080 ghcr.io/sergeyshevch/statuspage-exporter --statuspages=https://www.githubstatus.com, https://https://jira-software.status.atlassian.com
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
http_port: 8080
# Timeout for the http client
client_timeout: 2
# Count of retries for the http client
retry_count: 3
```

## Metrics Example

```
# HELP service_status Status of a service component, values 0 (operational) to 4 (major_outage)
# TYPE service_status gauge
service_status{component="API Requests",service="GitHub",status_page_url="https://www.githubstatus.com"} 1
service_status{component="Actions",service="GitHub",status_page_url="https://www.githubstatus.com"} 1
service_status{component="Codespaces",service="GitHub",status_page_url="https://www.githubstatus.com"} 1
service_status{component="Copilot",service="GitHub",status_page_url="https://www.githubstatus.com"} 1
service_status{component="Git Operations",service="GitHub",status_page_url="https://www.githubstatus.com"} 1
service_status{component="Issues",service="GitHub",status_page_url="https://www.githubstatus.com"} 1
service_status{component="Packages",service="GitHub",status_page_url="https://www.githubstatus.com"} 1
service_status{component="Pages",service="GitHub",status_page_url="https://www.githubstatus.com"} 1
service_status{component="Pull Requests",service="GitHub",status_page_url="https://www.githubstatus.com"} 1
service_status{component="Visit www.githubstatus.com for more information",service="GitHub",status_page_url="https://www.githubstatus.com"} 1
service_status{component="Webhooks",service="GitHub",status_page_url="https://www.githubstatus.com"} 1
# HELP service_status_fetch_duration_seconds Returns how long the service status fetch took to complete in seconds
# TYPE service_status_fetch_duration_seconds gauge
service_status_fetch_duration_seconds{status_page_url="https://www.githubstatus.com"} 7.0830795
```

## License Scan

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fsergeyshevch%2Fstatuspage-exporter.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fsergeyshevch%2Fstatuspage-exporter?ref=badge_large)



