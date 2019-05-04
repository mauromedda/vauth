FROM golang:alpine as builder
RUN apk --update add make git
RUN go get -u github.com/mauromedda/vauth
RUN /go/bin/vauth

FROM alpine
COPY --from=builder /go/bin/vauth /vauth

ENTRYPOINT [ "/vauth" ]
