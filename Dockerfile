FROM alpine

COPY ./pr0gramm-exporter /usr/bin/pr0gramm-exporter

ENTRYPOINT ["/usr/bin/pr0gramm-exporter"]
