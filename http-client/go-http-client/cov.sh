#!/bin/sh

out=./coverage.txt
oht=./coverage.html

go test \
	-race \
	-coverprofile="${out}" \
	-covermode=atomic \
	./...

go tool \
	cover \
	-html="${out}" \
	-o "${oht}"
