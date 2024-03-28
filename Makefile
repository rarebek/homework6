CURRENT_DIR=$(shell pwd)


build:
	CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o ${CURRENT_DIR}/bin/${APP} ${APP_CMD_DIR}/main.go

proto-p:
	protoc --go_out=genproto/product_service --go-grpc_out=genproto/product_service protos/product_service/product.proto

proto-u:
	protoc --go_out=genproto/user_service --go-grpc_out=genproto/user_service protos/user_service/user.proto

proto-m:
	protoc --go_out=genproto/message_service --go-grpc_out=genproto/message_service protos/message_service/message.proto

swag-gen:
	echo ${REGISTRY}
	swag init -g api/router.go -o api/docs