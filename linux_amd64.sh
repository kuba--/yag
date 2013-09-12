#!/bin/sh
GOOS=linux
GOARCH=amd64
PACKAGE=yag.tar.gz

mv "${GOOS}_${GOARCH}" "${GOOS}_${GOARCH}_$(date)"
mkdir -p "${GOOS}_${GOARCH}"
mv ${PACKAGE} "${GOOS}_${GOARCH}/" && cd "${GOOS}_${GOARCH}" && tar zxvf ${PACKAGE}

if [ -f "listener" ]; then
	pkill "[l]istener"
	echo "listener starting..."
	"./listener" > listener.log 2>&1 &
fi 

if [ -f "webserver" ]; then
	pkill "[w]ebserver"
	echo "webserver starting..."
	"./webserver" > webserver.log 2>&1 &
fi 

if [ -f "ttl" ]; then
	pkill "[t]tl"
	echo "ttl starting..."
	"./ttl" > ttl.log 2>&1 &
fi
