#!/bin/sh
echo "Deploy API Services..."

process_id=`ps -ef | grep ./trustcloudApiServices | grep -v "grep" | awk '{print $2}'`

# make sure LOG_DIR exists
if [ $process_id ]; then
    echo "Got running process...: " $process_id
	echo "Killing process with id ..." $process_id
	kill -9 $process_id
fi

cd /home/ubuntu

DATE=`date  +"%Y%m%d-%H%M"`
ARCHIVE_DIR=archive/trustcloud-api-services/$DATE

echo "Archiving current application to " $ARCHIVE_DIR
mkdir -p $ARCHIVE_DIR

mv trustcloud-api-services/* ./$ARCHIVE_DIR

cp -f /home/ubuntu/deploy/trustcloudApiServices /home/ubuntu/trustcloud-api-services
cp -rf /home/ubuntu/deploy/trustcloud/* /home/ubuntu/trustcloud-api-services
cp -rf /home/ubuntu/deploy/services /home/ubuntu/trustcloud-api-services

echo "Re-start API services server"

cd /home/ubuntu/trustcloud-api-services
nohup ./trustcloudApiServices >startup.log 2>&1 </dev/null &


process_id=`ps -ef | grep ./trustcloudApiServices | grep -v "grep" | awk '{print $2}'`


echo "API Services Server Deploy Completed..." $process_id

exit;