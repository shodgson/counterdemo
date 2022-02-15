#!/bin/bash

usage() { echo "Usage: $0 -n bucketname -p aws_profile -r aws_region" 1>&2; exit 1; }

while getopts n:p:r: flag
do
    case "${flag}" in
        n) bucketname=${OPTARG};;
        p) awsprofile=${OPTARG};;
        r) awsregion=${OPTARG};;
        *) usage ;;
    esac
done

if [ -z "${bucketname}" ] || [ -z "${awsprofile}" ] || [ -z "${awsregion}" ]; then
    usage
fi

echo "Good one"
echo "{
    \"Version\": \"2012-10-17\",
    \"Statement\": [
        {
            \"Sid\": \"PublicReadGetObject\",
            \"Effect\": \"Allow\",
            \"Principal\": \"*\",
            \"Action\": \"s3:GetObject\",
            \"Resource\": \"arn:aws:s3:::${bucketname}/*\"
        }
    ]
}" > /tmp/bucket_policy.json

aws s3api create-bucket --bucket ${bucketname} --region ${awsregion}  --create-bucket-configuration LocationConstraint=${awsregion} --profile ${awsprofile} \
  && aws s3api put-bucket-policy --bucket ${bucketname} --policy file:///tmp/bucket_policy.json --profile ${awsprofile} \
  && aws s3 website s3://${bucketname}/ --index-document index.html --error-document app/index.html --profile ${awsprofile}

