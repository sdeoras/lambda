#!/usr/bin/env bash

wget https://storage.googleapis.com/tensorflow/libtensorflow/libtensorflow-cpu-linux-x86_64-1.12.0.tar.gz
tar -zxvf libtensorflow-cpu-linux-x86_64-1.12.0.tar.gz
rm -rf libtensorflow-cpu-linux-x86_64-1.12.0.tar.gz include

cd myTFBin
GOOS=linux go build -o ../a.out .
