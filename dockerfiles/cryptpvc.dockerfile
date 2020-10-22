##################
### CI / BUILD ###
##################
FROM golang:1.14 as build

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

##########################
### CI / TEST / SANITY ###
##########################
FROM golang:1.14 as test-sanity
COPY . /workspace
WORKDIR /workspace

# Install csi-sanity
RUN curl -Lo csi-sanity.go https://raw.githubusercontent.com/kubernetes-csi/csi-test/master/cmd/csi-sanity/main.go &&\
    go build -o $GOPATH/bin/csi-sanity csi-sanity.go

# Wait for build and copy binary
COPY --from=build  /workspace/bin/cryptpvc-csi /workspace/bin/cryptpvc-csi

# Execute sanity test
CMD ["/workspace/hack/test-csi-sanity.sh"]

###############
### RUNTIME ###
###############
FROM alpine
LABEL maintainers="Daniel Sel <daniel-sel@hotmail.com>"
LABEL description="cryptpvc Kubernetes CSI Driver"

# Add util-linux to get a new version of losetup and.
RUN apk add --no-cache util-linux
COPY --from=build  /workspace/bin/cryptpvc-csi /cryptpvc-csi
ENTRYPOINT ["/cryptpvc-csi"]

# Skaffold debug autoconfig (this is a default setting and changes nothing)
ENV GOTRACEBACK=single 
