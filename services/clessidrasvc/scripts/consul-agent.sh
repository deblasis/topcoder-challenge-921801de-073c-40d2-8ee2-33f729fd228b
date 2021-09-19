#!/bin/bash

IP=$(ifconfig | sed -En 's/127.0.0.1//;s/.*inet (addr:)?(([0-9]*\.){3}[0-9]*).*/\2/p')
consul agent -bind=$IP -advertise=$IP -join=consul -node=$CONSUL_NODE -dns-port=53 -data-dir=/data -config-file=/consul/service.json -enable-local-script-checks