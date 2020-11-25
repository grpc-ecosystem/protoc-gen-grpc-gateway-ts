#!/bin/bash

source ./scripts/source.sh

CONF="karma.conf.js"

runTest ${1:-false} $CONF