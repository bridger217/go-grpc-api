download:
	@echo Download go.mod dependencies
	@go mod download

install-tools: download
	@echo Installing tools from tools.go
	@export GOBIN=$$PWD/bin && \
	export PATH=$$GOBIN:$$PATH && \
	cat pkg/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

user-service: install-tools
	@echo Compiling protos
	@export GOBIN=$$PWD/bin && \
	export PATH=$$GOBIN:$$PATH && \
	protoc --go_out=pkg/api/v1 --go_opt=paths=source_relative \
	--go-grpc_out=pkg/api/v1 --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out pkg/api/v1 --grpc-gateway_opt paths=source_relative \
	--proto_path=api/proto/v1 --proto_path=third_party \
	--openapiv2_out=docs/api/v1 --openapiv2_opt logtostderr=true \
	user_service.proto

clean:
	@rm -rf bin

server: user-service
	@go build -o bin/server cmd/server/main.go