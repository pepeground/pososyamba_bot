run_development:
	docker-compose run --rm bot bash

build: Dockerfile
	docker build --rm -f "Dockerfile" -t thesunwave/pososyamba:latest .

run:
	docker run --env-file .env.development pososyamba:latest
