#!/bin/bash
protoc -I greet --go_out=plugins=grpc:. greet/greetpb/greet.proto