build:
		@cd cmd && go guild -o ../bin/dockge-gitops

test:
		@go test ./...

generate:
		@go install github.com/matryer/moq@latest
		@go generate ./...

covtest:
		@mkdir -p "./bin"
		@go test -coverprofile=./bin/coverage.out ./...

covhtml: covtest
		@go tool cover -html=./bin/coverage.out