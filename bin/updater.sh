#!/bin/sh

BASEDIR=$(pwd)

cd ../src/trustcloud/batch
go run updater.go -bd=$BASEDIR/..
cd ../../../bin
