#!/bin/bash -x
set -ev

if [[ "$TRAVIS_OS_NAME" == "osx" ]]; then
    brew install -v llvm
    which llvm-link
    exit(0)
fi
