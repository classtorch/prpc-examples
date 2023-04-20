#!/bin/bash

protoc --proto_path=../third_party --go-prpc_out=:. --go_out=.  -I . api/user.proto
