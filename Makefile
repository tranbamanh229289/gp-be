BINARY_NAME=be
BUILD_DIR=bin
GQLGEN=github.com/99designs/gqlgen
DB_URL=postgres://postgres:postgres@localhost:5432/gp?sslmode=disable

#service
run: 
	@echo "🌐 Starting $(BINARY_NAME)..."
	go run ./cmd/api
build:
	@echo "🏗️ Building $(BINARY_NAME)..."
	go build -o ${BUILD_DIR}/${BINARY_NAME} ./cmd/api
clean:
	@echo "🧼 Cleaning ${BINARY_NAME} app..."
	rm -rf -o ${BUILD_DIR}

#docker
docker-infra: 
	docker compose -f docker-compose-infra.yml up -d

docker-dev:
	docker compose -f docker-compose-dev.yml up -d

docker-down:
	docker compose -f docker-compose-infra.yml down

#wire
wire: 
	wire gen ./internal/app

#graphql
gql-init:
	@echo "🚀 Initializing gqlgen.yml"
	go run ${GQLGEN} init
		
gql-gen: 
	@echo "🚀 Generating GraphQL code..."
	go run ${GQLGEN} generate

gql-clean:
	@echo "🧼 Cleaning generated files..."

#migrate
migrate-new: 
	@echo "🚀 Generating migration file  $(file).up.sql and ${file}.down.sql"
	migrate create -ext sql -dir ./internal/infrastructure/database/migration -seq $(file)

migrate-clean:
	@echo "🏗️ Migrating clean"
	migrate -path ./internal/infrastructure/database/migration -database ${DB_URL} drop -f
migrate-up:
	@echo "🏗️ Migrating up"
	migrate -path ./internal/infrastructure/database/migration -database ${DB_URL} -verbose up
	
migrate-down:
	@echo "🏗️ Migrating down"
	migrate -path ./internal/infrastructure/database/migration -database ${DB_URL} -verbose down

	