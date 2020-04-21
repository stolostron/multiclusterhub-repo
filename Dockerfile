# Build the manager binary
FROM golang:1.14 as builder

WORKDIR /workspace
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o start-repo ./cmd/repo

FROM registry.access.redhat.com/ubi8/ubi-minimal:latest
RUN microdnf install tar

WORKDIR /app
COPY --from=builder /workspace/start-repo .
COPY multiclusterhub/charts/ multiclusterhub/charts/
EXPOSE 3000
ENTRYPOINT /app/start-repo