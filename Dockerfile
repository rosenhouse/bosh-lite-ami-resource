FROM gliderlabs/alpine

ADD ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ADD bin /opt/resource
