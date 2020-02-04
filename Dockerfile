# Build the manager binary
FROM golang:1.13 as builder

WORKDIR /workspace
COPY main.go main.go
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o rhacm-repo main.go

FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

WORKDIR /app
COPY --from=builder /workspace/rhacm-repo .
COPY charts/ charts/
EXPOSE 3000
ENTRYPOINT /app/rhacm-repo