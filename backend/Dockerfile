FROM golang:1.23.3 as build

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/app /go/src/app/cmd/app/

FROM gcr.io/distroless/static-debian11
COPY --from=build /go/bin/app /
COPY config /config

EXPOSE 8080
CMD ["/app"]