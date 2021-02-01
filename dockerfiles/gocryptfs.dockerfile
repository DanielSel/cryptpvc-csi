### RUNTIME ###
FROM alpine
LABEL maintainers="Daniel Sel <daniel-sel@hotmail.com>"
LABEL description="cryptpvc Kubernetes CSI Driver"

# Add util-linux to get a new version of losetup.
RUN apk add --no-cache util-linux &&\
    apk add --no-cache -X http://dl-cdn.alpinelinux.org/alpine/edge/testing gocryptfs
ENTRYPOINT ["gocryptfs"]
