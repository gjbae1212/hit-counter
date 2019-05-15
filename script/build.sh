#!/bin/bash
set -e -o pipefail
trap '[ "$?" -eq 0 ] || echo "Error Line:<$LINENO> Error Function:<${FUNCNAME}>"' EXIT
cd `dirname $0` && cd ..
CURRENT=`pwd`

function start
{
   set_env
   local redis=`docker ps | grep redis | wc -l`
   if [ ${redis} -eq 0 ]
   then
     docker run --rm -d -p 6379:6379 --name redis redis:latest
   fi
   go build && ./hit-counter -tls=0  -addr=:8080
}

function test
{
   set_env
   go test -v $(go list ./... | grep -v vendor) --count 1 -race -coverprofile=$CURRENT/coverage.txt -covermode=atomic
}

function codecov
{
   /bin/bash <(curl -s https://codecov.io/bash)
}

function upload_config
{
   set_env
   gsutil cp $CURRENT/config/production.yaml $GCS_CONFIG_BUCKET/production.yaml
}


function set_env
{
   if [ -e $CURRENT/script/build_env.sh ]; then
     source $CURRENT/script/build_env.sh
   fi
}

CMD=$1
shift
$CMD $*
