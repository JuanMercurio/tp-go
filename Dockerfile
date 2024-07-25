FROM golang as build

WORKDIR /build

COPY . .

RUN go build -o app cmd/main.go
RUN rm .env
RUN mv .env-semideploy .env

CMD ["/build/app"]

# FROM scratch
# FROM alpine
#
# WORKDIR /deploy
#
# COPY --from=build /build/.env-deploy .env
# COPY --from=build /build/internal/adapters/config/apis_config.json api_config.json
# COPY --from=build /build/app .
#
# RUN chmod 777 app
#
# # analizar si usamos CMD
# CMD ["sh"]
