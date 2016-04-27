#!/bin/sh

BASEDIR=$(pwd)

cd ../src/trustcloud/batch
go run partner_api_check.go -bd=$BASEDIR/..
cd ../../../bin
