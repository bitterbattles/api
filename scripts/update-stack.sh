#!/bin/bash

NAME=$1
ARTIFACTDIR=../deployments
ARTIFACTNAME=$NAME.yml
S3BUCKET=bitterbattles-api-dev-cloudformation

echo Registering artifact...
aws s3 cp $ARTIFACTDIR/$ARTIFACTNAME s3://$S3BUCKET/$ARTIFACTNAME

echo Deploying artifact...
aws cloudformation update-stack \
    --stack-name $NAME \
    --template-url s3://$S3BUCKET/$ARTIFACTNAME