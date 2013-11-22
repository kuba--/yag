#!/bin/sh
GOOS=linux
GOARCH=amd64
PACKAGE=yag.tar.gz

set -x
mv "${GOOS}_${GOARCH}" "${GOOS}_${GOARCH}_$(date)"
mkdir -p "${GOOS}_${GOARCH}/logs"
mv ${PACKAGE} "${GOOS}_${GOARCH}/" && cd "${GOOS}_${GOARCH}" && tar zxvf ${PACKAGE}

for cmd in "listener" "webserver" "ttl"; do
	start-stop-daemon -v --stop --pidfile  "${HOME}/${GOOS}_${GOARCH}/${cmd}.pid"
	if [ $? -ne 0 ]
		then
		pkill "${cmd}"
	fi

	start-stop-daemon -v --start --pidfile "${HOME}/${GOOS}_${GOARCH}/${cmd}.pid" --background --exec "${HOME}/${GOOS}_${GOARCH}/${cmd}" -- -f="${HOME}/${GOOS}_${GOARCH}/config.json" -stderrthreshold=ERROR -log_dir="${HOME}/${GOOS}_${GOARCH}/logs"
done
