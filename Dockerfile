FROM alpine:latest

ADD  host_catalog /bin/host_catalog

ENTRYPOINT ["/bin/host_catalog"]