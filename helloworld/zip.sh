#!/usr/bin/env bash

cd bin && ./build.sh
cd ../
zip -r payload.zip hello.go go.mod bin/a.out bin/lib
