#!/usr/bin/env bash

proto_sources=$(find . -path "./api/fragma/*.proto")

protoc --go_out=. --go_opt=paths=source_relative ${proto_sources[@]}
