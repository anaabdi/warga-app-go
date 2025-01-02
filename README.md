# warga-app-go
a simple backend app that serve happiness for all warga

## Including:
- API: openapi spec with codegen to make life easier
- DB: postgres
- Cache: redis
- File Storage: minio (the S3 alternatives)

## How to Run
in vscode, 
- go to "Run and Debug" menu
- click Run button `Launch WargaAppGo`

### Config
run: `make config`

### Re-Generate API
run: `make gen-api`

### Test
#### Integration Test
run: `make integration-test`

#### Unit Test
run: `make unit-test`
