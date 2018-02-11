#!/bin/bash

git tag $1
GITHUB_TOKEN=`cat ./res/token.txt`
goreleaser --rm-dist
go install
