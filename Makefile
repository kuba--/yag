YAG=ubuntu@yag
PEM=${HOME}/.ssh/yag.pem
GOOS=darwin
GOARCH=amd64
PACKAGE=yag.tar
RUN_SCRIPT="./${GOOS}_${GOARCH}.sh"
CONFIG=config.json
LUA_SCRIPTS=add.lua get.lua ttl.lua

install:
	GOOS=${GOOS}	GOARCH=${GOARCH}	go install	./listener	./webserver	./ttl

deploy:	
	tar cvf	${PACKAGE}	-C ${GOPATH}/bin/${GOOS}_${GOARCH} listener webserver ttl
	tar rvf ${PACKAGE}	${CONFIG} ${LUA_SCRIPTS}
	gzip ${PACKAGE}

	scp -i ${PEM}	${PACKAGE}.gz	${YAG}:~/
	scp -i ${PEM}	${RUN_SCRIPT}	${YAG}:~/
	
all:	install	deploy	

clean:
	GOOS=${GOOS}	GOARCH=${GOARCH}	go clean -i -x ./listener ./webserver ./ttl
