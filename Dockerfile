FROM smarman/alp-go-base:latest AS build-env
ADD ./incident.go /go/src/damageReport/
RUN apk add --no-cache git && \
    cd /go/src/damageReport && \
    go get -u github.com/smarman85/dateRange && \
    go build incident.go

# final stage
FROM smarman/alpine-base
WORKDIR /app
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs
COPY --from=build-env /go/src/damageReport/incident /app
ENTRYPOINT ./incident
