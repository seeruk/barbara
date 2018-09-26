#!/usr/bin/env bash

set -e

scriptDir=$(dirname $0)

pushd "$scriptDir/.." > /dev/null || exit 1

    qtdeploy build desktop ./cmd/barbara

    exec cmd/barbara/deploy/linux/barbara

popd > /dev/null || exit 1
