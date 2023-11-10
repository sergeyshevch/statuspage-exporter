FROM golang:1.21 as builder

COPY . /build/
WORKDIR /build

RUN CGO_ENABLED=0 go build -ldflags "-s -w" . 

# Now copy it into our base image.
FROM gcr.io/distroless/static-debian12:nonroot

COPY --chown=nonroot:nonroot --from=builder /build/statuspage-exporter /usr/local/bin/statuspage-exporter
USER nonroot
ENTRYPOINT ["/usr/local/bin/statuspage-exporter"]