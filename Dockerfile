FROM golang:1.21 as build

WORKDIR /go/secure-share

COPY go.mod go.sum /go/secure-share/
RUN go mod download && go mod verify

COPY . .

RUN --mount=type=cache,id=gomod,target=/go/pkg/mod \
    --mount=type=cache,id=gobuild,target=/root/.cache/go-build \
    GOOS=linux CGO_ENABLED=0 go build -v -o /go/bin/secure-share /go/secure-share/cmd/server


FROM alpine:latest
COPY --from=build /go/bin/secure-share /
COPY web /web
CMD ["/secure-share"]
