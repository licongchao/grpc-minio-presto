########################
# STEP 1 build executable binary
########################
FROM golang:alpine AS builder

ENV GOPROXY https://goproxy.cn,direct
####################################################
# Security
####################################################
ENV USER=appuser
ENV UID=10001

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

# Change apk 
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
# RUN echo -e http://mirrors.ustc.edu.cn/alpine/v3.16/main/ > /etc/apk/repositories
# Install git
RUN apk update && apk add --no-cache git

COPY . /app
WORKDIR /app

# Fetch dependencies
RUN go get -d -v

# Build the binary
###################
# RUN go build -o /go/bin/main
#######################
# Removing debug information & 
# disabling cross compilation
#######################
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/main

#########################
# STEP 2 build smaller image
#########################
FROM scratch

#############################################
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

COPY --from=builder /go/bin/main /app

# Use an unprivileged user
USER appuser:appuser
EXPOSE 10000

ENTRYPOINT ["/app","serve"]