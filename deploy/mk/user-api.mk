VERSION = latest

SERVER_NAME= user
SERVER_TYPE= api

DOCKER_REPO_TEST = registry.cn-hangzhou.aliyuns.com/easy-chat/${SERVER_NAME}-${SERVER_TYPE}-dev

VERSION_TEST = $(VERSION)

APP_NAME_TEST = easy-chat-${SERVER_NAME}-${SERVER_TYPE}-dev

DOCKER_FILE_TEST = ./deploy/dockerfile/Dockerfile_${SERVER_NAME}_${SERVER_TYPE}_dev

build-test:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/${SERVER_NAME}-${SERVER_TYPE} ./apps/${SERVER_NAME}/${SERVER_TYPE}/${SERVER_NAME}.go
	docker build . -f ${DOCKER_FILE_TEST} --no-cache -t ${APP_NAME_TEST}

tag-test:
	@echo 'create tag ${VERSION_TEST}'
	docker tag ${APP_NAME_TEST} ${DOCKER_REPO_TEST}:${VERSION_TEST}

publish-test:
	@echo 'publish ${VERSION_TEST} to ${DOCKER_REPO_TEST}'
	docker push ${DOCKER_REPO_TEST}:${VERSION_TEST}

release-test: build-test tag-test publish-test