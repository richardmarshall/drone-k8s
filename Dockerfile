FROM alpine:3.7

LABEL maintainer="Richard Marshall <rgm@linux.com>"

ENV KUBECTL_VERSION="v1.9.3"

RUN apk add --no-cache ca-certificates \
    && apk add --no-cache --virtual .deps curl \
    && curl -L https://storage.googleapis.com/kubernetes-release/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl -o /usr/local/bin/kubectl \
    && chmod +x /usr/local/bin/kubectl \
    && apk del --purge .deps

ADD drone-k8s /bin/

ENTRYPOINT ["/bin/drone-k8s"]
