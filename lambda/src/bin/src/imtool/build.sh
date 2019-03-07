#!/usr/bin/env bash

# install for whatever OS/arch you are on
go install

# a.out needs to be a linux/amd64 binary since it runs within a
# cloud function environment. Since it wraps Tensorflow C library
# it used cgo for build and can therefore not be cross-compiled.
GOOS=linux go build -o a.out .
