FROM golang:1.8 as build-env

WORKDIR /go/src/app
COPY . .

RUN go-wrapper download 
RUN CGO_ENABLED=0 GOOS=linux go-wrapper install -a -ldflags '-extldflags "-static"' 

CMD ["go-wrapper", "run"] # ["app"]

FROM docker:stable

COPY --from=build-env /go/bin/app /app
CMD ["/app"]

