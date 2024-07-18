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


.DEFAULT_GOAL := vet
.PHONY := fmt run build