#!/bin/sh

BASEDIR=$(pwd)

cd ../src/trustcloud/batch
go run batch_runner.go -bd=$BASEDIR/..
cd ../../../bin
