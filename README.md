# docker-cloud-aws-sync
[![Build Status](https://travis-ci.org/PoweredLocal/docker-cloud-aws-sync.svg?branch=master)](https://travis-ci.org/PoweredLocal/docker-cloud-aws-sync)
[![Code Climate](https://codeclimate.com/github/PoweredLocal/docker-cloud-aws-sync/badges/gpa.svg)](https://codeclimate.com/github/PoweredLocal/docker-cloud-aws-sync)

[![Docker Hub](http://dockeri.co/image/pwred/docker-cloud-aws-sync)](https://hub.docker.com/r/pwred/docker-cloud-aws-sync/)

Synchronizing Docker Cloud node IPs with an AWS security group

## Overview

A simple Golang script that synchronizes your Docker Cloud node list with specified AWS security group - 
so that your Docker Cloud nodes can always access AWS services (eg. RDS, Memcached, etc)

If all of your resources are in AWS, this is not an issue – you can keep everything inside the same VPC, but what if some of your nodes are hosted elsewhere but need to access some of the AWS services?

## Requirements

- Go >= 1.6
- Official AWS Go SDK
- Official Docker Cloud Go SDK

## How it works

docker-cloud-aws-sync will fetch current node list, flush AWS security group and re-add all nodes during startup.
Afterwards, this script will keep listening to Docker Cloud events and if there are node updates – it will update AWS
security accordingly.

Security groups are part of EC2 service so don't forget to attach a policy with corresponding EC2 permissions (roles) to your API key.

## Running it as a Docker container

There is a Dockerfile included and there is a [public repository](https://hub.docker.com/r/pwred/docker-cloud-aws-sync/) in Docker Hub.

So either ```docker build ./``` or use the public image ```docker run -d -e DOCKER_CLOUD_xxx=xxx -e AWS_xxx=xxx pwred/docker-cloud-aws-sync```

## Environment variables

Script relies on several environment variables to access your Docker Cloud and AWS:

```bash
DOCKER_CLOUD_USER 
# required - Your username in Docker Cloud

DOCKER_CLOUD_KEY
# required - API key

DOCKER_CLOUD_NAMESPACE 
# (optional) - If necessary, organization namespace

DOCKER_CLOUD_TAG 
# (optional) - If specified, only nodes that have this tag will be processed

AWS_ACCESS_KEY_ID 
# (required) - AWS key id

AWS_SECRET_ACCESS_KEY 
# (required) - AWS key secret

AWS_SG_ID 
# (required) - AWS security group id where rules should be pushed to

AWS_REGION 
# (required) - AWS region (AWS will use default if not set)

```

## To do

- Parse live Docker Cloud event stream
- Add node tag support

## NB

Please note that EC2 security groups allow up to 50 rules per group. Therefore, if you have more than 50 nodes you have to create multiple security groups and add Docker Cloud nodes to them based on tags. Or you have to migrate everything to a single cloud provider! :)

## License

The MIT License (MIT)
Copyright (c) 2016 PoweredLocal

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
