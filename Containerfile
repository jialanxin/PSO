FROM docker.io/library/golang:1.25-bookworm

WORKDIR /src

ENV CGO_ENABLED=0 \
    GOFLAGS=-mod=readonly \
    GOTOOLCHAIN=local

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["go", "test", "./..."]
