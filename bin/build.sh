cd ../src
pushd trustcloud
GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o ../../out/trustcloudApiServices
popd
cd ../bin
