FROM golang:alpine AS build-env
RUN mkdir /go/src/app && apk update && apk add git
ADD . /go/src/app/
WORKDIR /go/src/app
RUN go mod download && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .

FROM scratch
WORKDIR /app
COPY --from=build-env /go/src/app/main .
EXPOSE 5003
ENTRYPOINT [ "./main" ]
