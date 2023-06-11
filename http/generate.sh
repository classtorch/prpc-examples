#!/bin/bash

protoc --proto_path=../third_party --go-prpc_opt http_generate_grpc=false --go-prpc_out=:. --go_out=.  -I . api/user.proto
