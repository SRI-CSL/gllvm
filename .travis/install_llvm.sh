#!/bin/bash -x
set -ev

if [[ "$TRAVIS_OS_NAME" == "osx" ]]; then
    brew install -v llvm
fi
