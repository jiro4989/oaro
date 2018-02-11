#!/bin/bash

git tag $1
export GITHUB_TOKEN=`cat ./res/token.txt`
goreleaser --rm-dist
go install
