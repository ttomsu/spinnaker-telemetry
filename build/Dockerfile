FROM gcr.io/cloud-builders/go as builder

WORKDIR /workspace

RUN pwd

RUN ls -a

ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -v ./

# -------

FROM mirror.gcr.io/library/alpine
RUN apk add --no-cache ca-certificates

COPY --from=builder /workspace/server /server

CMD ["/server"]
