FROM golang as build

WORKDIR /build

COPY . .

RUN go build -o app cmd/main.go

FROM scratch

WORKDIR /app

COPY --from=build /build/.env-deploy /app/.env
COPY --from=build /build/internal/adapters/config/apis_config.json /app/api_config.json
COPY --from=build /build/app /app

# analizar si usamos CMD
ENTRYPOINT ["/app"]
