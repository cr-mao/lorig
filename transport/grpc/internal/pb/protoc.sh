#!/bin/bash

## protoc --gofast_out=.. --go-grpc_out=.. *.proto

protoc --proto_path=.  --go_out=.  --go-grpc_out=.    *.proto
