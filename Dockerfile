FROM gcr.io/distroless/static-debian11:nonroot

COPY "./statuspage-exporter" /usr/local/bin/statuspage-exporter

ENTRYPOINT ["statuspage-exporter"]
