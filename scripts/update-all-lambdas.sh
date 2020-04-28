#!/bin/bash

for dir in ../cmd/*
do
    cmd=${dir##*/}
    sh update-lambda.sh $cmd
done