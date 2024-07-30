
.DEFAULT_GOAL:= vet
.PHONY: fmt run build mock dev air

DIRBIN=bin/
BIN=$(DIRBIN)app
MAIN=cmd/main.go
DEPENDENCIES=

build: fmt
	go build -o $(BIN) $(MAIN)

# formatea segun el estandar de go
fmt:
	go fmt ./...

# avisa de posibles errores que no se ve en el compilador como:
# printf("%s") o codigo inancanzable
vet: build
	go vet ./...

run:
	go run $(MAIN)

clean:
	rm -fr $(DIRBIN)

mock:
	go generate ./...

test: 
	go test ./...

dev:
	~/go/bin/air -c .air.toml

cover:
	go test ./... -cover -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
