# Statuspage Exporter
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fsergeyshevch%2Fstatuspage-exporter.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fsergeyshevch%2Fstatuspage-exporter?ref=badge_shield)
![Build Status](https://github.com/sergeyshevch/statuspage-exporter/workflows/CI/badge.svg)
[![License](https://img.shields.io/github/license/sergeyshevch/statuspage-exporter)](/LICENSE)
[![Release](https://img.shields.io/github/release/sergeyshevch/statuspage-exporter.svg)](https://github.com/golangci/golangci-lint/releases/latest)
[![Docker](https://img.shields.io/docker/pulls/sergeyshevch/statuspage-exporter)](https://hub.docker.com/r/golangci/golangci-lint)

Statuspage exporter exports metrics from given statuspages as prometheus metrics.

## Supported statuspage engines:
- Statuspage.io (Widely used statuspage engine. For example by [GitHub](https://www.githubstatus.com)). You can check that statuspage is supported by this engine by checking that it has a [/api/v2/components.json](https://www.githubstatus.com/api/v2/components.json) endpoint.

## Some popular statuspages:
- [GitHub](https://www.githubstatus.com)
- [Atlassian (Jira/Confluence/etc)](https://status.atlassian.com/)

## Running exporter

You can run the exporter with docker, kubernetes, or just as a binary.

### Docker
```bash
docker run -p 8080:8080 ghcr.io/sergeyshevch/statuspage-exporter --statuspage=https://www.githubstatus.com
```

### kubernetes
TODO: Provide helm chart

### Binary
TODO: Provide docs for binary running

## Configuration

You can provide configuration using configuration file or environment variables.

Configuration file must be named as '.statuspage-exporter.yaml' and should be placed in the home directory or same directory as the binary.

Environment variable names are the same as configuration file keys but in upper case and with underscores instead of dots.

### Configuration file example
```yaml
http_port: 8080
# Delay between requests to the statuspages
fetch_delay: 5
# Timeout for the http client
client_timeout: 2
# List of the targets to scrape
statuspages:
  - https://githubstatus.com
```

## Metrics Example
```
# HELP service_status Status of a service component, values 0 (operational) to 4 (major_outage)
# TYPE service_status gauge
service_status{component="API Requests",service="GitHub",status_page_url="https://status.github.com"} 1
service_status{component="Actions",service="GitHub",status_page_url="https://status.github.com"} 1
service_status{component="Codespaces",service="GitHub",status_page_url="https://status.github.com"} 1
service_status{component="Copilot",service="GitHub",status_page_url="https://status.github.com"} 1
service_status{component="Git Operations",service="GitHub",status_page_url="https://status.github.com"} 1
service_status{component="Issues",service="GitHub",status_page_url="https://status.github.com"} 1
service_status{component="Packages",service="GitHub",status_page_url="https://status.github.com"} 1
service_status{component="Pages",service="GitHub",status_page_url="https://status.github.com"} 1
service_status{component="Pull Requests",service="GitHub",status_page_url="https://status.github.com"} 1
service_status{component="Visit www.githubstatus.com for more information",service="GitHub",status_page_url="https://status.github.com"} 1
service_status{component="Webhooks",service="GitHub",status_page_url="https://status.github.com"} 1
# HELP statuspage_exporter_build_info A metric with a constant '1' value labeled by version, revision, branch, and goversion from which statuspage_exporter was built.
# TYPE statuspage_exporter_build_info gauge
statuspage_exporter_build_info{branch="",goversion="go1.19.2",revision="",version=""} 1
```

## License Scan

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fsergeyshevch%2Fstatuspage-exporter.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fsergeyshevch%2Fstatuspage-exporter?ref=badge_large)



