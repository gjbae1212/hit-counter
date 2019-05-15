#!/bin/bash

set -e -o pipefail
trap '[ "$?" -eq 0 ] || echo "Error Line:<$LINENO> Error Function:<${FUNCNAME}>"' EXIT
cd `dirname $0` && cd ..
CURRENT=`pwd`

function deploy
{
   echo "yes" | gcloud app deploy production.yaml --promote
}

function test
{
  go test -v $(go list ./... | grep -v vendor) --count 1 -race
}

function download_config
{
  BUCKET_PATH=`cat /root/config/bucket_path`
  gsutil cp $BUCKET_PATH $CURRENT/production.yaml
}

CMD=$1
shift
$CMD $*