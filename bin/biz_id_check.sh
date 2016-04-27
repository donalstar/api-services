#!/bin/sh

BASEDIR=$(pwd)

cd ../src/trustcloud/batch
go run id_check_runner.go -bd=$BASEDIR/..
cd ../../../bin
