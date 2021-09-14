# gcloudx
extra features for accessing the Google Cloud Platform

## requirements

Authenticated GCP user account with relevant permissions.

    gcloud auth application-default login

## install

    go install github.com/emicklei/gcloudx/cmd/gcloudx@latest

## pubsub

    gcloudx pubsub publish -h

    NAME:
    gcloudx pubsub publish - publish a document from file

    USAGE:
    gcloudx pubsub publish [command options] [arguments...]

    OPTIONS:
    -p value  GCP project identifier
    -t value  PubSub topic identifier (short name)
    -f value  File containing the payload

## iam

    gcloudx iam roles -h   

    NAME:
    gcloudx iam roles - list all permissions assigned to a member

    USAGE:
    gcloudx iam roles [arguments...]

&copy; 2021 ernestmicklei.com MIT License.    