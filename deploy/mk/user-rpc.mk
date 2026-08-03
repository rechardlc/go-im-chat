VERSION=latest

# 定义服务器名称和类型
SERVER_NAME=user
SERVER_TYPE=rpc

# 定义镜像名称和标签
VERSION_TEST=$(VERSION)
APP_NAME_TEST=easy-im-${SERVER_NAME}-${SERVER_TYPE}-test
DOCKER_FILE_TEST=./deploy/dockerfile/Dockerfile.${SERVER_NAME}_${SERVER_TYPE}_dev

VERSION_PROD=$(VERSION)
APP_NAME_PROD=easy-im-${SERVER_NAME}-${SERVER_TYPE}-prod
DOCKER_FILE_PROD=./deploy/dockerfile/Dockerfile.${SERVER_NAME}_${SERVER_TYPE}_prod

build-test:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/${SERVER_NAME}-${SERVER_TYPE}-test ./apps/${SERVER_NAME}/${SERVER_TYPE}/${SERVER_NAME}.go
	docker build -t ${APP_NAME_TEST}:${VERSION_TEST} -f ${DOCKER_FILE_TEST} .

run-test:
	docker run -d --name ${APP_NAME_TEST} ${APP_NAME_TEST}:${VERSION_TEST}