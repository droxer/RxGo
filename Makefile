default: deps test

vet:
	go vet ./...

deps:
	go get -d -v ./...
	go list -f '{{range .TestImports}}{{.}} {{end}}' ./... | xargs -n1 go get -d

test:
	go test ./...

.PHONY: deps vet test