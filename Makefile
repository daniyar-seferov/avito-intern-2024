DOCKER_DIR=${CURDIR}/build/dev
DOCKER_YML=${DOCKER_DIR}/docker-compose.yml
ENV_NAME="avito-tender"
APP_NAME="avitoapp"

-include ./build/dev/.env

.PHONY: compose-up
compose-up:
	docker-compose -p ${ENV_NAME} -f ${DOCKER_YML} up

.PHONY: compose-down
compose-down: ## terminate local env
	docker-compose -p ${ENV_NAME} -f ${DOCKER_YML} stop

.PHONY: compose-rm
compose-rm: ## remove local env
	docker-compose -p ${ENV_NAME} -f ${DOCKER_YML} rm -fvs

.PHONY: compose-rs
compose-rs: ## remove previously and start new local env
	make compose-rm
	make compose-up


# WITH PROD ENVIRONMENT
.PHONY: docker-build
docker-build:
	docker build . -t ${APP_NAME}

.PHONY: docker-run
docker-run:
	docker run --rm -p 8080:8080 ${APP_NAME}

.PHONY: docker-rs
docker-rs:
	make docker-build
	make docker-run