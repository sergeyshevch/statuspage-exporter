FROM distroless

COPY "./statuspage-exporter" /usr/local/bin/statuspage-exporter

ENTRYPOINT ["statuspage-exporter"]
