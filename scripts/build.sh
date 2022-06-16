#!/bin/bash

if [ -d output ]; then rm output/*; fi
if [[ ! -d output ]]; then mkdir output; fi
go build -o output/server ./server/server.go
npm run build --prefix demo