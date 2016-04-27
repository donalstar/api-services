#!/bin/sh

BASEDIR=$(pwd)

cd ../src/trustcloud/batch
go run fc_company_check_runner.go -bd=$BASEDIR/..
cd ../../../bin
