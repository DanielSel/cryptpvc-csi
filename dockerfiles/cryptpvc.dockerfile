### BUILD / CI ###
FROM golang:1.14 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy relevant sources from repo
COPY cmd/cryptpvc-csi/ cmd/cryptpvc-csi/
COPY pkg/ pkg/
COPY hack/ hack/

## Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o ./bin/cryptpvc-csi ./cmd/cryptpvc-csi


### RUNTIME ###
FROM alpine
LABEL maintainers="Daniel Sel <daniel-sel@hotmail.com>"
LABEL description="cryptpvc Kubernetes CSI Driver"

# Add util-linux to get a new version of losetup and.
RUN apk add --no-cache util-linux
COPY --from=builder  /workspace/bin/cryptpvc-csi /cryptpvc-csi
ENTRYPOINT ["/cryptpvc-csi"]

# Skaffold debug autoconfig (this is a default setting and changes nothing)
ENV GOTRACEBACK=single 
