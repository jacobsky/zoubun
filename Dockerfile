FROM --platform=$BUILDPLATFORM golang:1.25 AS builder

WORKDIR /build

ARG TARGETOS TARGETARCH

COPY go.mod go.sum ./

RUN go mod download
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

COPY . .
# run code generation prior to complilation
RUN sqlc generate
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 go build -o server /build/cmd/server/

# Prod layer for actual proper deployments.
# NOTE: It would usually be better to use a more secure runtime such as one from chainguard,
# for this example, it doesn't specifically make a difference (and I lack a subscription/access to the direct repos anyways)
FROM scratch AS prod
WORKDIR /app
COPY --from=builder /build/server ./server

EXPOSE 3000
CMD ["/app/server"]
