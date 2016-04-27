#!/bin/sh

BASEDIR=$(pwd)

cd ../src/trustcloud/batch
go run indiv_id_check_runner.go -bd=$BASEDIR/..
cd ../../../bin
