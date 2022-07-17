APP_NAME=autonotif-scheduler
LOAD_DIR=./app-images
DOCKER_NETWORK=autonotif-network

.PHONE: docker-network
docker-network:
	docker network create $(DOCKER_NETWORK)

.PHONE: docker-postgres
docker-postgres:
	docker run -d -p 5432:5432 \
		--net $(DOCKER_NETWORK) \
		--name autonotif-postgres \
		--env-file db/autonotif-postgres.env \
		postgres:14.4-alpine

.PHONE: docker-psql
docker-psql:
	docker run -it --rm --network $(DOCKER_NETWORK) postgres:14.4-alpine psql -h autonotif-postgres autonotif -U autonotif_usr

.PHONY: migrate
migrate:
	docker run --rm -v $(shell pwd)/db/migrations:/migrations --network $(DOCKER_NETWORK) migrate/migrate -path=/migrations/ -database $(url) up

.PHONY: rollback
rollback:
	docker run --rm -v $(shell pwd)/db/migrations:/migrations --network $(DOCKER_NETWORK) migrate/migrate -path=/migrations/ -database $(url) down 1

.PHONE: docker-load
docker-load:
	docker load -i $(LOAD_DIR)/$(APP_NAME)-latest.tar

.PHONY: docker-run
docker-run:
	docker run --rm -d \
		-p 8080:8080 \
		-v $(shell pwd)/config:/app/config \
		-e CONFIG_PATH=./config/config.yaml \
		--net $(DOCKER_NETWORK) \
		--name $(APP_NAME) \
		$(APP_NAME):latest

.PHONE: run
run: docker-load docker-run

.PHONY: stop
stop:
	docker stop $(APP_NAME)
