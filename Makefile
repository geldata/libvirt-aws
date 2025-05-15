build:
	go build -o bin/ ./...
	chmod +x bin/*

tidy:
	go mod tidy

fmt:
	gofmt -s -l -w .

test:
	go test -v ./...

vet:
	go vet ./...