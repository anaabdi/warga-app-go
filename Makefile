.PHONY: install-api-codegen
install-api-codegen:
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

.PHONY: gen-api
gen-api:
	oapi-codegen -package=api -generate "types,spec,chi" api/v1/v1.yaml > api/v1/api.gen.go

.PHONY: config
config:
	cp .env.example .env

# located in deployments/local/docker-compose.yaml
.PHONY: infra-up
infra-up:
	docker-compose -f deployments/local/docker-compose.yaml up -d

.PHONY: infra-down
infra-down:
	docker-compose -f deployments/local/docker-compose.yaml down && \
	docker-compose -f deployments/local/docker-compose.yaml down --volumes

.PHONY: infra-logs
infra-logs:
	docker-compose -f deployments/local/docker-compose.yaml logs -f

.PHONY: infra-ps
infra-ps:
	docker-compose -f deployments/local/docker-compose.yaml ps
	

.PHONY: unit-test
unit-test:
	go test -v ./... -v -coverprofile=coverage.out

.PHONY: unit-test-coverage
unit-test-coverage:
	go tool cover -html=coverage.out

.PHONY: integration-test
integration-test:
	go test -v ./... -tags=integration