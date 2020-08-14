#!/bin/sh

if [ ! -x /app/goip_exporter ]; then
  chmod 755 /app/goip_exporter
fi

cd /app

if [ -z "$CONFIG_FILE" ]
then
    ./goip_exporter --goip.address=${GOIP_ADDRESS:-192.168.1.189}  --goip.user=${GOIP_USER:-admin} --goip.password=${GOIP_PASSWORD:-admin} ${GOIP_EXTRA_PARAMS}
else
    ./goip_exporter -config-file $CONFIG_FILE
fi
