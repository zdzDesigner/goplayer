#!/usr/bin/env bash

arg=$1
ENV=${arg#*=}
# echo "${ENV^^}"

if [[ $1 ]]; then
export __ENV__="${ENV^^}"
fi


go run .
# go run -race .
