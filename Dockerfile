#build stage
FROM golang:1.23.6 AS build-env

ADD . /build_dir
WORKDIR /build_dir

RUN rm go.sum && go mod tidy

WORKDIR /build_dir/cmd
RUN CGO_ENABLED=0 go build -o /quiz_backend_core

#Final stage
FROM alpine:3.16.0

EXPOSE 80

WORKDIR /
COPY --from=build-env /quiz_backend_core /

ENTRYPOINT ["/quiz_backend_core"]
