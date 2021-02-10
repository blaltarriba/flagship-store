install:
	go get github.com/gorilla/mux
	go get github.com/google/uuid
	go get github.com/stretchr/testify/assert
	go get github.com/stretchr/testify/mock

build: ## Build project
	docker run --rm -it -v "$$GOPATH":/gopath -v "$$(pwd)":/app -e "GOPATH=/gopath" -w /app golang:1.15.7 sh -c 'CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags="-s" -o flagship-store'

	docker build -t flagship-store .

run: ## Start project container
	docker run --rm -it -p 3080:3080 flagship-store

test:
	go test ./... -v