FROM golang:1.10.0-stretch as builder
WORKDIR /go/src/github.com/bitrise-team/addons-template-service
COPY . .
RUN mkdir -p bin
RUN CGO_ENABLED=0 go build  -o ./bin/addons-template-service -a -ldflags '-extldflags "-static"' .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/github.com/bitrise-team/addons-template-service/bin/addons-template-service /usr/local/bin/addons-template-service

CMD ["addons-template-service"]
