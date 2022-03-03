#!/bin/bash
set -e -o pipefail
trap '[ "$?" -eq 0 ] || echo "Error Line:<$LINENO> Error Function:<${FUNCNAME}>"' EXIT
cd `dirname $0` && cd ..
CURRENT=`pwd`

function deploy
{
   test
   make_wasm "production"
   upload_deployment_config
   gsutil cp $GCS_CONFIG_BUCKET/production.yaml $CURRENT/production.yaml
   rm $CURRENT/cloudbuild.yaml || true
   echo "yes" | gcloud app deploy production.yaml --promote
   git checkout cloudbuild.yaml
   rm  $CURRENT/production.yaml || true
}

function make_wasm
{
  local phase=$1
  echo "WASM PHASE=$phase"

  # copy wasm_exec.js
  cp $(go env GOROOT)/misc/wasm/wasm_exec.js $CURRENT/public/

  GOOS=js GOARCH=wasm go build -ldflags="-s -w -X main.phase=$phase" -o $CURRENT/view/hits.wasm $CURRENT/wasm/main.go
  gzip $CURRENT/view/hits.wasm
  mv $CURRENT/view/hits.wasm.gz $CURRENT/view/hits.wasm
}


function upload_deployment_config
{
   set_env
   gsutil cp $CURRENT/script/production.yaml $GCS_CONFIG_BUCKET/production.yaml
}

function run
{
   # make_wasm "local"
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
   go test -v $(go list ./... | grep -v vendor | grep -v wasm) --count 1 -timeout 120s
}

function set_env
{
   ulimit -n 1000
   if [ -e $CURRENT/script/build_env.sh ]; then
     source $CURRENT/script/build_env.sh
   fi
}

CMD=$1
shift
$CMD $*
