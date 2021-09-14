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


## bq

    gcloudx bq deps -h  

    NAME:
    gcloudx bq deps - bq deps PROJECT(.|:)DATASET.VIEW,...

    USAGE:
    gcloudx bq deps [command options] [arguments...]

    OPTIONS:
    -o value  output file with DOT notation (default: "bigquery.dot")

### example

    gcloudx bq deps -o g.dot myproject-id.my_dataset.my_view && cat g.dot | dot -Tpng > gcloudx-deps-bigquery.png && open gcloudx-deps-bigquery.png

&copy; 2021 ernestmicklei.com MIT License.    