#!/bin/bash -x
set -ev

if [[ "$TRAVIS_OS_NAME" == "osx" ]]; then
    wget http://csl.sri.com/users/iam/llvm_lite-6.0.0.high_sierra.bottle.1.tar.gz
    brew install -v ./llvm_lite-6.0.0.high_sierra.bottle.1.tar.gz
    #brew install -v llvm
fi
