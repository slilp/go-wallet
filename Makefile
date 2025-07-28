run-dev:
	source .env && air -c .air.toml

test:
	go test -v ./...

tidy:
	go mod tidy
	
generate-gen:
	@rm -rf ./internal/servers/api_gen/*
	@mkdir -p ./internal/servers/api_gen
	$(shell go env GOPATH)/bin/oapi-codegen -package=api_gen --generate types -o ./internal/servers/api_gen/openapi_types.gen.go ./docs/server.yml 
	$(shell go env GOPATH)/bin/oapi-codegen -package=api_gen --generate gin -o ./internal/servers/api_gen/openapi_api.gen.go ./docs/server.yml

generate-mock:
	go generate ./...

