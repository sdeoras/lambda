#!/usr/bin/env bash
set -ex
protoc -I . api.proto --go_out=plugins=grpc:.
