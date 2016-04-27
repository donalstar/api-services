#!/bin/sh

BASEDIR=$(pwd)

cd ../src/trustcloud/batch
go run background_checker.go -bd=$BASEDIR/..
cd ../../../bin
