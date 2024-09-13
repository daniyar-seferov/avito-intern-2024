BIN_DIR=${CURDIR}/bin
DOCKER_DIR=${CURDIR}/build/dev
DOCKER_YML=${DOCKER_DIR}/docker-compose.yml
ENV_NAME="avito-tender"
APP_NAME="avitoapp"
MIGRATE_DIR=${CURDIR}/migrations

# -include ./build/dev/.env
-include ./build/prod/.env

.PHONY: compose-up
compose-up:
	docker-compose -p ${ENV_NAME} -f ${DOCKER_YML} up --force-recreate --build

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

# MIGRATION
.PHONY: install-goose
install-goose:
	GOBIN=${BIN_DIR} go install github.com/pressly/goose/v3/cmd/goose@latest

.PHONY: goose-migrate-up
goose-migrate-up:
	${BIN_DIR}/goose -dir ${MIGRATE_DIR} postgres "\
		user=${POSTGRES_USERNAME}\
		dbname=${POSTGRES_DATABASE}\
		password=${POSTGRES_PASSWORD}\
		port=${POSTGRES_PORT}\
		sslmode=disable" up

.PHONY: goose-migrate-down
goose-migrate-down:
	${BIN_DIR}/goose -dir ${MIGRATE_DIR} postgres "\
		user=${POSTGRES_USERNAME}\
		dbname=${POSTGRES_DATABASE}\
		password=${POSTGRES_PASSWORD}\
		port=${POSTGRES_PORT}\
		sslmode=disable" down

.PHONY: goose-migrate-prod-up
goose-migrate-prod-up:
	${BIN_DIR}/goose -dir ${MIGRATE_DIR} postgres "\
		user=${POSTGRES_USERNAME}\
		dbname=${POSTGRES_DATABASE}\
		password=${POSTGRES_PASSWORD}\
		port=${POSTGRES_PORT}\
		host=${POSTGRES_HOST}\
		sslmode=disable" up