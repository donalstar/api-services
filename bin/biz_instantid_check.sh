#!/bin/sh

BASEDIR=$(pwd)

cd ../src/trustcloud/batch
go run biz_instantid_check_runner.go -bd=$BASEDIR/..
cd ../../../bin
