#!/bin/bash
set -ev

OPWD=$PWD
NAME=widgets_playground


qtdeploy -docker build windows_64_static && rm -rf rcc* moc* vendor && git clean -f && git reset --hard && docker rmi therecipe/qt:windows_64_static
cd $OPWD/deploy/windows && zip -9qrXy ../${NAME}_windows_amd64.zip * && cd $OPWD && rm -rf $OPWD/deploy/windows 

qtdeploy -docker build linux_static && rm -rf rcc* moc* vendor && git clean -f && git reset --hard && docker rmi therecipe/qt:linux_static
cd $OPWD/deploy/linux && zip -9qrXy ../${NAME}_linux_amd64.zip * && cd $OPWD && rm -rf $OPWD/deploy/linux

cd $(go env GOPATH)/src/github.com/therecipe/qt/internal/docker/darwin && ./build_static.sh && cd $OPWD
qtdeploy -docker build darwin_static && rm -rf rcc* moc* vendor && git clean -f && git reset --hard && docker rmi therecipe/qt:darwin_static
cd $OPWD/deploy/darwin && zip -9qrXy ../${NAME}_darwin_amd64.zip * && cd $OPWD && rm -rf $OPWD/deploy/darwin

qtdeploy -docker build js && rm -rf rcc* moc* vendor && git clean -f && git reset --hard && docker rmi therecipe/qt:js
cd $OPWD/deploy/js && zip -9qrXy ../${NAME}_js.zip * && cd $OPWD && rm -rf $OPWD/deploy/js 

