run-dev:
	docker compose -f docker-compose.dev.yml up

test:
	go test -v ./...

tidy:
	go mod tidy
	
generate-gen:
	@rm -rf ./internal/restapis/api_gen/*
	@mkdir -p ./internal/restapis/api_gen
	$(shell go env GOPATH)/bin/oapi-codegen -package=api_gen --generate types -o ./internal/port/restapis/api_gen/openapi_types.gen.go ./docs/server.yml 
	$(shell go env GOPATH)/bin/oapi-codegen -package=api_gen --generate gin -o ./internal/port/restapis/api_gen/openapi_api.gen.go ./docs/server.yml

generate-mock:
	go generate ./...

