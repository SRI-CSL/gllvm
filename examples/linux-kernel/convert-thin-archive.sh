#!/bin/bash
for lib in $*;
    do ar -t $lib |xargs ar rvs $lib.new;
done