FROM golang:1.23.0 AS build
WORKDIR /src
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /src/tfidf ./cmd/api/main.go

FROM alpine:3.21.3 AS production
COPY --from=build /src/tfidf /bin/tfidf
COPY --from=build /src/web/templates /web/templates


EXPOSE 8080
CMD ["./bin/tfidf"]