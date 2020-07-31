ARG work_dir=/go/src/github.com/Bnei-Baruch/feed-api

FROM golang:1.14-alpine3.11 as build

LABEL maintainer="edoshor@gmail.com"

ARG work_dir

ENV GOOS=linux \
	CGO_ENABLED=0

RUN apk update && \
    apk add --no-cache \
    git

WORKDIR ${work_dir}
COPY . .

RUN go test -v $(go list ./...) \
    && go build


FROM alpine:3.11
ARG work_dir
WORKDIR /app
COPY --from=build ${work_dir}/feed-api .

EXPOSE 8080
CMD ["./feed-api", "server"]
