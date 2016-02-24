#!/bin/bash

protoc -I. pkg/crowdsound/crowdsound_service.proto --go_out=plugins=grpc:.
protoc -I. pkg/playsource/playsource_service.proto --go_out=plugins=grpc:.

