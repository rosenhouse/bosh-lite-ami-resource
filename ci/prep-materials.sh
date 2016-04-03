#!/bin/bash

set -e -x -u

export GOPATH=/go
OUT=$PWD/image-materials

cp bosh-lite-ami-resource/Dockerfile $OUT/

mkdir $OUT/certs
cp /etc/ssl/certs/ca-certificates.crt $OUT/certs/

mkdir -p $GOPATH/src/github.com/rosenhouse
cp -R bosh-lite-ami-resource $GOPATH/src/github.com/rosenhouse/

mkdir $OUT/bin
for action in check in out; do
  go build -o $OUT/bin/$action github.com/rosenhouse/bosh-lite-ami-resource/actions/$action
done

cd $OUT
pwd
find .
