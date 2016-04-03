FROM gliderlabs/alpine

ADD certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ADD bin /opt/resource
