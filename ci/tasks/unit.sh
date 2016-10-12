#!/usr/bin/env bash

set -e -x
source voyager-cli/ci/tasks/util.sh

check_param GITHUB_USER
check_param GITHUB_PASSWORD

echo -e "machine github.com\n  login $GITHUB_USER\n  password $GITHUB_PASSWORD" >> ~/.netrc

export GOPATH=$PWD
export PATH=$PATH:$GOPATH/bin
mkdir -p $GOPATH/src/github.com/RackHD/
cp -r voyager-cli $GOPATH/src/github.com/RackHD/voyager-cli

pushd $GOPATH/src/github.com/RackHD/voyager-cli
  make deps
  make build
  make unit-test
  echo "Unit test complete."
popd
