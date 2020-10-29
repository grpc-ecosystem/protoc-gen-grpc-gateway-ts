#!/bin/bash

go run ./ &
pid=$!

./node_modules/.bin/karma start karma.conf.ci.js
pkill -P $pid
