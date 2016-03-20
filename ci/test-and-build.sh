#!/bin/bash

set -e -x -u

pwd
ls
export GOPATH=$PWD/go
cd $GOPATH/src/github.com/rosenhouse/bosh-lite-ami-resource

go install ./actions/check
go install ./actions/in
go install ./actions/out
