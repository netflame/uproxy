#!/usr/bin/env sh

APPNAME=uproxy

go build -o ${APPNAME}_${GOOS}_${GOARCH} \
  && ln -s ${APPNAME}_${GOOS}_${GOARCH} ${APPNAME}
