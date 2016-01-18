#!/bin/bash

protoc -I. pkg/crowdsound/crowdsound_service.proto --go_out=plugins=grpc:.

