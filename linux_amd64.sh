#!/bin/sh
GOOS=linux
GOARCH=amd64
PACKAGE=yag.tar.gz
SCRIPTS=add.lua get.lua ttl.lua

mv "${GOOS}_${GOARCH}" "${GOOS}_${GOARCH}_$(date)"
mkdir -p "${GOOS}_${GOARCH}"
mv ${PACKAGE} "${GOOS}_${GOARCH}/" && cd "${GOOS}_${GOARCH}" && tar zxvf ${PACKAGE} && cd ${HOME}


if [ -f "${GOOS}_${GOARCH}/listener" ]; then
	pkill "[l]istener"
	echo "listener starting..."
	"${GOOS}_${GOARCH}/listener" -f ./${GOOS}_${GOARCH} > listener.log 2>&1 &
fi 

if [ -f "${GOOS}_${GOARCH}/webserver" ]; then
	pkill "[w]ebserver"
	echo "webserver starting..."
	"${GOOS}_${GOARCH}/webserver" -f ./${GOOS}_${GOARCH} > webserver.log 2>&1 &
fi 

if [ -f "${GOOS}_${GOARCH}/ttl" ]; then
	pkill "[t]tl"
	echo "ttl starting..."
	"${GOOS}_${GOARCH}/ttl" -f ./${GOOS}_${GOARCH} > ttl.log 2>&1 &
fi
