#!/usr/bin/env bash
BASE_DIR=$(dirname $(dirname $(readlink -f "$0")))
# deploy app name
APP_STRING=$1
DEPLOY_APPS=(${APP_STRING//,/ })

function usage() {
  echo "Deploy app"
  echo ""
  echo "$0 <app array> <method>"
  echo ""
  echo "e,g: deploy.sh web"
  echo ""
}

function checkOptApp() {
  # 检查部署应用
  ALL_APPS="web"
  for APP in ${DEPLOY_APPS[@]}
  do
    if [[ $ALL_APPS =~ (^|[[:space:]])$APP($|[[:space:]]) ]];then
      echo "[OK] app name: $APP"
    else
      usage
      echo "Please specify valid app name, e,g: $ALL_APPS"
      exit 1
    fi
  done
}

function checkOpts() {
  checkOptApp
}

function buildApps() {
  for APP in ${DEPLOY_APPS[@]}
  do
    echo ""
    echo "===> Building $APP..."
    echo ""
    rm -fr $APP
    go build -o $APP cmd/$APP/main.go
    if [ $? -ne 0 ]; then
      echo "Error building $APP"
      exit 1
    fi
  done
}

#############
# main
#############
checkOpts
buildApps


