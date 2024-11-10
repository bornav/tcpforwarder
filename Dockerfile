# Start by building the application.
# FROM golang:1.20-alpine as build

FROM golang:1.23.1 AS build
WORKDIR /usr/src/tcpforwarder
COPY . .

# RUN go test ./...
RUN CGO_ENABLED=0 GO111MODULE=on go build -a -installsuffix nocgo -o ./tcpforwarder ./

# # Now copy it into our base image.
# FROM gcr.io/distroless/static-debian11:nonroot
# FROM scratch
FROM  ubuntu:20.04
# FROM golang:1.19
COPY --from=build /usr/src/tcpforwarder/tcpforwarder /usr/bin/tcpforwarder
# COPY --from=build /usr/bin/bash /bin/bash
COPY startup.sh /startup.sh

# ENTRYPOINT [ "/usr/bin/tcpforwarder" ]
ENTRYPOINT [ "/startup.sh" ]
# CMD []

# LABEL org.opencontainers.image.title tcpforwarder
# LABEL org.opencontainers.image.description "Forward tcp packets from the source to the destination"
# LABEL org.opencontainers.image.licenses MIT
