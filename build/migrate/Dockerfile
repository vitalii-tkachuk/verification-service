FROM golang:1.19 as base

WORKDIR /go/src/github.com/vitalii-tkachuk/verification-service

ENV CGO_ENABLED=0

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

FROM base AS build

COPY . .

ARG SKAFFOLD_GO_GCFLAGS

RUN --mount=target=. \
    --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=cache,target=/root/.cache/go-build/ \
    go build -gcflags="${SKAFFOLD_GO_GCFLAGS}" -trimpath -o /go/bin/migrate cmd/migrate/main.go

FROM scratch
COPY --from=build /go/bin/migrate /go/bin/migrate
COPY --from=build /go/src/github.com/vitalii-tkachuk/verification-service/migrations /migrations

ENTRYPOINT ["/go/bin/migrate"]
