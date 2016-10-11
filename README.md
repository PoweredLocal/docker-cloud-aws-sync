# docker-cloud-aws-sync
[![Build Status](https://travis-ci.org/dusterio/docker-cloud-aws-sync.svg?branch=master)](https://travis-ci.org/dusterio/docker-cloud-aws-sync)

Synchronizing Docker Cloud nodes with AWS security groups

## Overview

A simple Golang script that synchronizes your Docker Cloud nodes with specified AWS security group - 
so that your Docker Cloud nodes can always access AWS services (eg. RDS)

If all of your resources are in AWS, this is not an issue – you can keep everything inside the same VPC, but what if some of
your nodes are hosted elsewhere, but need to access some of the AWS nodes or services?

## Requirements

- Go >= 1.6
- Official AWS Go SDK
- Official Docker Cloud Go SDK

## How it works

docker-cloud-aws-sync will fetch current node list, flush AWS security group and re-add all nodes during startup.
Afterwards, this script will keep listening to Docker Cloud events and if there are node updates – it will update AWS
security accordingly.

## Running it as a Docker container

There is a Dockerfile included and there is a [public repository](https://hub.docker.com/r/dusterio/docker-cloud-aws-sync/) in Docker Hub.

So either ```docker build ./``` or use the public image ```docker run -e DOCKER_CLOUD_xxx dusterio/docker-cloud-aws-sync```

## Environment variables

Script relies on several environment variables to access your Docker Cloud and AWS:

```
DOCKER_CLOUD_USER (required) - Your username in Docker Cloud
DOCKER_CLOUD_KEY (required) - API key
DOCKER_CLOUD_NAMESPACE (optional) - If necessary, organization namespace

AWS_KEY (required) - AWS key id
AWS_SECRET (required) - AWS key secret

```