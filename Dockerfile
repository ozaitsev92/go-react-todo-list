FROM golang:latest AS builder

COPY . /app
WORKDIR /app
RUN ls

RUN CGO_ENABLED=0 go build -ldflags '-s -w -extldflags "-static"' -o /app/appbin ./cmd/apiserver

FROM debian:stable-slim
LABEL MAINTAINER Author ozaitsev92@gmail.com

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        ca-certificates  \
        netbase \
    && rm -rf /var/lib/apt/lists/ \
    && apt-get autoremove -y && apt-get autoclean -y

RUN adduser --home "/appuser" --disabled-password appuser --gecos "appuser,-,-,-"
USER appuser

RUN mkdir -p /appuser/app

COPY --from=builder /app/appbin /appuser/app
COPY configs /appuser/app/configs
COPY static /appuser/app/static

WORKDIR /appuser/app

EXPOSE 8000

CMD ["./appbin"]