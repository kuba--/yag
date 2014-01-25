#!/bin/sh
GOOS=linux
GOARCH=amd64
PACKAGE=yag.tar.gz

set -x
mv "${GOOS}_${GOARCH}" "${GOOS}_${GOARCH}_$(date)"
mkdir -p "${GOOS}_${GOARCH}/logs"
mv ${PACKAGE} "${GOOS}_${GOARCH}/" && cd "${GOOS}_${GOARCH}" && tar zxvf ${PACKAGE}

for cmd in "listener" "webserver" "ttl"; do
	EXEC="${HOME}/${GOOS}_${GOARCH}/${cmd} -- -f=${HOME}/${GOOS}_${GOARCH}/config.json -stderrthreshold=WARNING -log_dir=${HOME}/${GOOS}_${GOARCH}/logs"

	start-stop-daemon --stop --pidfile "${HOME}/${cmd}.pid"
	if [ $? -ne 0 ]
		then
		pkill "${cmd}"
	fi
	start-stop-daemon --start --background --make-pid --pidfile "${HOME}/${cmd}.pid" --exec $EXEC
done
