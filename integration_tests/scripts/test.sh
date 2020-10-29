#!/bin/bash

go run ./ &
pid=$!

./node_modules/.bin/karma start karma.conf.js
pkill -P $pid
