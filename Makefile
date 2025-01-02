.PHONY: install-api-codegen
install-api-codegen:
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

.PHONY: gen-api
gen-api:
	oapi-codegen -package=api -generate "types,spec,chi" api/v1/v1.yaml > api/v1/api.gen.go

.PHONY: config
config:
	cp .env.example .env