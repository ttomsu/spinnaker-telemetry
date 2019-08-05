FROM mirror.gcr.io/library/golang as builder

ADD . /workspace
WORKDIR /workspace

ENV GOPROXY=https://proxy.golang.org
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux

RUN go install -mod=vendor cmd/server/main.go

RUN ls -a /go/bin

# -------

FROM mirror.gcr.io/library/alpine
RUN apk add --no-cache ca-certificates

COPY --from=builder /go/bin/main /server

ENTRYPOINT ["/server"]
