#!/bin/bash

LAMBDANAME=$1
LAMBDADIR=../cmd/$LAMBDANAME
OUTPUTDIR=../temp
ARTIFACTNAME=$LAMBDANAME.zip
S3BUCKET=bitterbattles-api-dev-lambda

echo Running tests...
go test $LAMBDADIR

echo Building code...
mkdir -p $OUTPUTDIR
rm -f $OUTPUTDIR/$LAMBDANAME*
GOOS=linux GOARCH=amd64 go build -o $OUTPUTDIR/$LAMBDANAME $LAMBDADIR
echo ok $OUTPUTDIR/$LAMBDANAME

echo Packaging artifact...
zip -j $OUTPUTDIR/$ARTIFACTNAME $OUTPUTDIR/$LAMBDANAME
rm $OUTPUTDIR/$LAMBDANAME

echo Registering artifact...
aws s3 cp $OUTPUTDIR/$ARTIFACTNAME s3://$S3BUCKET/$ARTIFACTNAME

echo Deploying artifact...
aws lambda update-function-code \
    --function-name $LAMBDANAME \
    --s3-bucket $S3BUCKET \
    --s3-key $ARTIFACTNAME \
    --publish \
    --no-dry-run

echo Cleaning up...
rm -f $OUTPUTDIR/$LAMBDANAME*