#!/usr/bin/env bash

rm -rf ./target/fserv ./target/templates ./target/static
go build
mkdir -p target
mv ./fserv ./target
cp -r ./env.sh ./templates ./static ./target
bash ./minify.sh ./target/static/css
