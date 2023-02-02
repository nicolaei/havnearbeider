#!/usr/bin/env sh

# I haven't implemented any safe-guards, so this is how we do it!

sudo rm -rf build/ filesystems images/alpine-goat.tar
go build
