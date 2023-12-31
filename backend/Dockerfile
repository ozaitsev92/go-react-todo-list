FROM golang:latest as build

WORKDIR /go/src/app
COPY . .
RUN ls -la
RUN pwd

RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/app /go/src/app/cmd/apiserver/

FROM gcr.io/distroless/static-debian11
COPY --from=build /go/bin/app /
COPY configs /configs

EXPOSE 8080
CMD ["/app"]