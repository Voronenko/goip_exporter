FROM alpine:latest as alpine
RUN apk --update add ca-certificates

FROM scratch
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

EXPOSE 9177

COPY goip_exporter /goip_exporter

ENTRYPOINT [ "/goip_exporter" ]
