#!/bin/bash
source ./scripts/source.sh

CONF="karma.conf.ci.js"

runTest false $CONF
runTest true $CONF

