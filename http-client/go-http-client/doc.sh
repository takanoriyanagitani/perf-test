#!/bin/sh

port=6060

godoc \
  -v \
  -play \
  -timestamps \
  -http localhost:${port}
