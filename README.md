# gcloudx
extra features for accessing the Google Cloud Platform

## requirements

Authenticated GCP user account with relevant permissions.

    gcloud auth application-default login

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

&copy; 2021 ernestmicklei.com MIT License.    