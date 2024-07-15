FROM alpine:latest as build

RUN apk add --no-cache go

ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH

WORKDIR /app

# despues ver que copiar para performance
COPY . .

RUN go build -o app cmd/main.go

FROM scratch

COPY --from=build /app/cmd/app /app

# analizar si usamos CMD
ENTRYPOINT ["/app"]