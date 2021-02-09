install:
	go get github.com/gorilla/mux
	go get github.com/google/uuid
	go get github.com/stretchr/testify/assert

build: ## Build project
	docker run --rm -it -v "$$GOPATH":/gopath -v "$$(pwd)":/app -e "GOPATH=/gopath" -w /app golang:1.15.7 sh -c 'CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags="-s" -o flagship-store'

	docker build -t flagship-store .

env-start: ## Start project container
	docker run --rm -it -p 3080:10000 flagship-store

test:
	go test ./... -v