#!/bin/bash

set -e -x

[[ -s "/home/emc/.gvm/scripts/gvm" ]] >/dev/null 2>/dev/null
source "/home/emc/.gvm/scripts/gvm" >/dev/null 2>/dev/null

PROJECT_PATH=$GOPATH/src/github.com/RackHD

cleanUp()
{
  # Don't exit on error here. All commands in this cleanUp must run,
  #   even if some of them fail
  set +e

  # Clean up all cloned repos
  cd $GOPATH
  rm -rf $GOPATH/src
}

trap cleanUp EXIT

pushd $PROJECT_PATH/voyager-cli
  echo "Testing Mission Control Center"

  make deps
  make build
  make integration-test

  echo "Mission Control Center PASS\n\n"
popd

exit
