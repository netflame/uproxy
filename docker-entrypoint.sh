#!/usr/bin/env sh
cmd=${1}
nohup ${cmd} 1>>out.log 2>>err.log \
    && tail -f out.log err.log