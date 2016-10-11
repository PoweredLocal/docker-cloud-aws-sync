# docker-cloud-aws-sync
Synchronizing Docker Cloud nodes with AWS security groups

## Overview

A simple Golang script that synchronizes your Docker Cloud nodes with specified AWS security group - 
so that your Docker Cloud nodes can always access AWS services (eg. RDS)

If all of your resources are in AWS, this is not an issue – you can keep everything inside the same VPC, but what if some of
your nodes are hosted elsewhere, but need to access some of the AWS nodes or services?

## How it works

docker-cloud-aws-sync will fetch current node list, flush AWS security group and re-add all nodes during startup.
Afterwards, this script will keep listening to Docker Cloud events and if there are node updates – it will update AWS
security accordingly.
