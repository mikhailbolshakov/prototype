# current version
VERSION ?= 0.1.0
# docker registrty url
DOCKER_URL = docker.medzdrav.ru

.PHONY: build proto

proto: ## Generates proxies for proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./proto/users/*.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./proto/tasks/*.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./proto/chat/*.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./proto/services/*.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./proto/bp/*.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./proto/config/*.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./proto/sessions/*.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./proto/webrtc/*.proto

build: proto ## builds the main
	go build -o bin/ cmd/main.go

build-local: proto ## builds all services locally
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./config/bin/main ./config/cmd/main.go
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./users/bin/main ./users/cmd/main.go

docker-build-local: build-local  ## Build the docker images for all services
	@echo Building images
	docker build . -f ./config/Dockerfile.local -t $(DOCKER_URL)/prototype/config:$(VERSION) --no-cache
	docker build . -f ./users/Dockerfile.local -t $(DOCKER_URL)/prototype/users:$(VERSION) --no-cache

docker-push-local: docker-build-local ## Build and push docker images to the repository
	@echo Pushing images
	docker push $(DOCKER_URL)/prototype/config:$(VERSION)
	docker push $(DOCKER_URL)/prototype/users:$(VERSION)

docker-build:  ## Build the docker images for all services (build inside)
	@echo Building images
	docker build . -f ./config/Dockerfile -t $(DOCKER_URL)/prototype/config:$(VERSION) --no-cache
	docker build . -f ./users/Dockerfile -t $(DOCKER_URL)/prototype/users:$(VERSION) --no-cache

docker-push: docker-build ## Build and push docker images to the repository
	@echo Pushing images
	docker push $(DOCKER_URL)/prototype/config:$(VERSION)
	docker push $(DOCKER_URL)/prototype/users:$(VERSION)

docker-build-all:  ## Build the docker images for all services
	@echo Building images
	docker build . -f ./config/Dockerfile -t $(DOCKER_URL)/prototype/config:$(VERSION) --no-cache
	docker build . -f ./users/Dockerfile -t $(DOCKER_URL)/prototype/users:$(VERSION) --no-cache

docker-push-all: docker-build-all ## Build and push docker images to the repository
	@echo Pushing images
	docker push $(DOCKER_URL)/prototype/config:$(VERSION)
	docker push $(DOCKER_URL)/prototype/users:$(VERSION)

