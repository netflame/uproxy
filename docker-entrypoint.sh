#!/usr/bin/env sh
nohup ./uproxy 1>>out.log 2>>err.log &
tail -f out.log err.log
