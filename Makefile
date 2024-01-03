build:
		@cd cmd && go guild -o ../bin/dockge-gitops

test:
		@go test ./...

generate:
		@go install github.com/matryer/moq@latest
		@go generate ./...