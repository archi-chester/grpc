#!/bin/bash
protoc --go_out=plugins=grpc:. blogpb/blog.proto