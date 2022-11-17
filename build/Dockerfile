FROM golang:1.19 as build

WORKDIR /go/src/github.com/vitalii-tkachuk/verification-service

COPY . .

# todo will be able to speed up build if cache dependencies in separate step
ARG SKAFFOLD_GO_GCFLAGS
RUN CGO_ENABLED=0 go build -gcflags="${SKAFFOLD_GO_GCFLAGS}" -trimpath -o /go/bin/verification-service cmd/api/main.go

FROM scratch
COPY --from=build /go/bin/verification-service /go/bin/verification-service

ENTRYPOINT ["/go/bin/verification-service"]
