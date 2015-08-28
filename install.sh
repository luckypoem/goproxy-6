#!/bin/sh

export GOPATH=`pwd`

go install local
go install remote

echo 'Ok!!'
