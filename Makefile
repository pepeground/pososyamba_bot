DOCKER_COMPOSE_RUN=docker-compose run bot

build_image: Dockerfile
	docker build --rm -f "Dockerfile" -t thesunwave/pososyamba:latest .

run: build_image
	docker run --env-file .env.development pososyamba:latest

clean:
	rm -f ./main

build_container:
	docker-compose build

build_devel: build_container
	${DOCKER_COMPOSE_RUN} "go build /application/cmd/pososyamba_bot/main.go"
	chmod a+x ./main

up: build_devel
	docker-compose up

down:
	docker-compose down

sh:
	${DOCKER_COMPOSE_RUN} /bin/bash -l
