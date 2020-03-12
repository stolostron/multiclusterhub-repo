# Build the manager binary
FROM golang:1.13 as builder

WORKDIR /workspace
COPY main.go main.go
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o start-repo main.go

FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

WORKDIR /app
COPY --from=builder /workspace/start-repo .
COPY multiclusterhub/charts/ multiclusterhub/charts/
EXPOSE 3000
ENTRYPOINT /app/start-repo