#!/usr/bin/env bash

rm -rf target
go build
mkdir target
mv ./fserv ./target
cp -r ./env.sh ./templates ./static ./target
bash ./minify.sh ./target/static/css