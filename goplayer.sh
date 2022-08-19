#!/bin/bash

LISTEN_FILE=.fsnotify


if [[ $1 == "-next" ]]; then
  echo "next" > $LISTEN_FILE
  exit
fi

if [[ $1 == "-prev" ]]; then
  echo "prev" > $LISTEN_FILE
  exit
fi

if [[ $1 == "-pause" ]]; then
  echo "pause" > $LISTEN_FILE
  exit
fi
