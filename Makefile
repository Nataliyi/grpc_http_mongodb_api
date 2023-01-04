PRJ = project
APP = grpc-api
BIN = $(APP)
VER = v1

ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
PWD:= 

.PHONY: info ### App info
info:
	make -v
	sudo docker version --format 'Client: {{ .Client.Version}} Server: {{ .Server.Version }}'
	go version
	echo "namespace:"$(PRJ) "appname:"$(APP) "binary-name:"$(BIN) "version:"$(VER)

.PHONY: build ### Build app
build:
	CGO_ENABLED=0 go build -o $(BIN) -v

.PHONY: build-image ### Build images
build-image:
	docker volume create --name init-db --opt type=none --opt device=$(ROOT_DIR)/mongo-init.js --opt o=bind
	env | grep PWD && \
	sudo -E bash -c 'sudo docker-compose build'

up: ### Up docker-compose
	env | grep PWD && \
	sudo -E bash -c 'docker-compose up -d' && docker-compose logs -f grpc-api
.PHONY: up

down: ### Down docker-compose
	sudo docker-compose down 
.PHONY: down

