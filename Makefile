BINARY_NAME=be
BUILD_DIR=bin
GQLGEN=github.com/99designs/gqlgen

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
	docker compose -f ./build/docker-compose-infra.yml up -d

docker-dev:
	docker compose -f ./build/docker-compose-dev.yml up -d
	
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