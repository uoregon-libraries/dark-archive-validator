#!/usr/bin/env bash
set -eu

success=1
make
gb test rules -run "ExampleEngine.*" >test.log 2>&1 || success=0
if [[ $success == 1 ]]; then
  rm test.log
  echo "Success"
  exit 0
fi

go run testsplit.go test.log || cat test.log
for gotfile in $(find . -name "*.got"); do
  wantfile=${gotfile%.got}.want
  diff -u $wantfile $gotfile || true
done

rm -f test.log test.*.want test.*.got
