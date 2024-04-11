FROM golang:1.21.6-bullseye

SHELL ["/bin/bash", "-c"]

RUN apt-get update && apt-get -y upgrade \
    && apt-get autoremove -y \
    && rm -rf /var/lib/apt/lists/* \
    && apt-get -y clean

WORKDIR /build

COPY go.mod *.go ./
COPY plugins/*.go ./plugins/
COPY plugins/wasm ./plugins/wasm/
COPY pluginapi/* ./pluginapi/
RUN go mod tidy \
    && go mod download \
    && go build -o /extproc

FROM golang:1.21.6-bullseye

COPY --from=0 /extproc /extproc

EXPOSE 50051

ARG EXAMPLE=noop
ENTRYPOINT [ "/extproc" ]
CMD [ "noop" ]
