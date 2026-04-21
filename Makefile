# Regenerate restapi server code from swagger.yaml (requires github.com/go-swagger/go-swagger).
# Install: go install github.com/go-swagger/go-swagger/cmd/swagger@latest
SWAGGER ?= $(shell go env GOPATH)/bin/swagger

.PHONY: swagger-gen
swagger-gen:
	$(SWAGGER) generate server --target=. --name=RedbudWayAPI --spec=./swagger.yaml --principal=interface{}

.PHONY: build
build:
	go build ./...
