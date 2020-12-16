UID=$(shell id -u)
USRENAME=`id -un`
USER_PARAM=--user $(shell id -u):$(shell id -g)
APP_VOLUME_PARAM=-v `pwd`/.:/application
APP_PORT_PARAM=-p 80:3000/tcp
CONTAINER_NAME=thesunwave/pososyamba:latest
DOCKER_RUN=docker run -it ${USER_PARAM} ${APP_VOLUME_PARAM} ${APP_PORT_PARAM} ${CONTAINER_NAME}
DOCKER_COMPOSE_RUN=docker-compose run  ${USER_PARAM}  bot

build_prod: Dockerfile
	docker build --rm -f "Dockerfile" -t thesunwave/pososyamba:latest .

run_prod: build_image
	${DOCKER_RUN}--env-file .env.development

clean:
	rm -f ./main

build_devel: clean
	${DOCKER_COMPOSE_RUN} /bin/sh -c "go build -work /application/cmd/pososyamba_bot/main.go"
	chmod a+x ./main

up: build_devel
	docker-compose up

down:
	docker-compose down

sh:
	${DOCKER_COMPOSE_RUN} /bin/bash -l
