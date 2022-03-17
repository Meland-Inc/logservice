FROM golang:1.17.5 as builder

ARG platform
ARG mainpath
ARG COMMITID
ARG BUILDENV

RUN mkdir /app

WORKDIR /go/src/app

# copy 源码到容器
ADD ./ ./

# # ## build main
RUN GOOS="${platform}" GOARCH=amd64 CGO_ENABLED=1 go build \
    -tags "build_bindata ${BUILDTAGS}" \
    -ldflags '-extldflags "-static"  -X bian.CommitId=${COMMITID}  -X bian.BuildEnv=${BUILDENV}'  \
    -o "/app" \
    "./main.go"

FROM alpine

RUN apk --no-cache add  \
    ca-certificates \
    libc6-compat \
    libstdc++ \
    file \
    tzdata

ARG env

RUN mkdir /opt/app

COPY --from=builder /app/* /opt/app

RUN chmod +x /opt/app

WORKDIR /opt/app

VOLUME /opt/app/logs

EXPOSE 50001

ENTRYPOINT ["/opt/app/main"]