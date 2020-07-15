# Build the manager binary
FROM golang:1.13 as builder

WORKDIR /workspace
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o start-repo ./cmd/repo

FROM registry.access.redhat.com/ubi8/ubi-minimal:latest
RUN microdnf install tar
RUN microdnf update

LABEL org.label-schema.vendor="Red Hat" \
      org.label-schema.name="multiclusterhub-repo" \
      org.label-schema.description="Helm repo that serves charts for the Red Hat Advanced Cluster Management installer" \
      org.label-schema.license="Red Hat Advanced Cluster Management for Kubernetes EULA"

WORKDIR /app
COPY --from=builder /workspace/start-repo .
COPY multiclusterhub/charts/ multiclusterhub/charts/
EXPOSE 3000
RUN chmod -R 777 /app/multiclusterhub/charts
ENTRYPOINT /app/start-repo