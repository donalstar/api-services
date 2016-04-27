#!/bin/sh
echo "Deploy API Services\n"

ENV=$1
if [ -z "$ENV" ]; then
  ENV="dev"
fi

if [ "$ENV" != "dev" ] && [ "$ENV" != "prod" ]; then
    echo "USAGE: deploy.sh [dev|prod]"
fi

echo "\tDeploying to $ENV..."

if [ "$ENV" == "dev" ]; then
    DEST_SERVER=ec2-54-235-249-103.compute-1.amazonaws.com
    KEY_FILE=trust_cloud_dev.pem
else
#    DEST_SERVER=ec2-107-20-240-138.compute-1.amazonaws.com     #prod server

    DEST_SERVER=ec2-184-72-216-160.compute-1.amazonaws.com    #dev4 server


#    KEY_FILE=trust_cloud.pem
    KEY_FILE=trust_cloud_dev.pem
fi

KEYS_DIR=/Users/donal/TrustCloud/keys

ARCHIVE=api-services.tar

./build.sh

mkdir -p ./trustcloud/config
cp ../src/trustcloud/config/trustcloud.gcfg trustcloud/config
cp ../src/trustcloud/config/$ENV.gcfg trustcloud/config/env.gcfg

mkdir -p ./trustcloud/public
cp ../resources/* trustcloud/public

tar cf $ARCHIVE ./execute.sh -C ../out trustcloudApiServices -C ../src trustcloud/templates trustcloud/cert -C ../bin trustcloud services

scp -i $KEYS_DIR/$KEY_FILE $ARCHIVE ubuntu@$DEST_SERVER:/tmp

ssh -i $KEYS_DIR/$KEY_FILE ubuntu@$DEST_SERVER << EOT

cd /home/ubuntu
mkdir -p deploy
cd deploy
rm -rf *
tar  xvf /tmp/api-services.tar
./execute.sh

EOT

rm -rf tmp
rm $ARCHIVE

rm -rf ./trustcloud

echo "Deploy completed"


exit;

