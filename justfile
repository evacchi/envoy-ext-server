
image_name := "envoy-extproc-sdk-go-examples"
image_tag := `git rev-parse HEAD`

default:
    just --list

update:
    go get -u ./ && just tidy

tidy:
    go mod tidy

format:
    go fmt *.go

unit-test: 
    echo "TBD"

integration-test: 
    echo "TBD"

coverage: 
    echo "TBD"

run example="noop" *flags="":
    go run *.go {{flags}} {{example}}

build *flags="":
    go build {{flags}}
    cd plugins/data    && go build -buildmode=plugin
    cd plugins/dedup   && go build -buildmode=plugin
    cd plugins/digest  && go build -buildmode=plugin
    cd plugins/echo    && go build -buildmode=plugin
    cd plugins/masker  && go build -buildmode=plugin
    cd plugins/noop    && go build -buildmode=plugin
    cd plugins/timer   && go build -buildmode=plugin
    cd plugins/trivial && go build -buildmode=plugin
    cd plugins/wasm    && go build -buildmode=plugin

clean:
    rm plugins/**/*.so

containerize tag=image_tag *flags="": 
    docker build . -t {{image_name}}:{{tag}} {{flags}}
#    [[ -d extproc ]] && rm -rf extproc || true
#    mkdir -p extproc \
#        && cp ../*.go extproc/ \
#        && cp ../go.mod ../go.sum extproc/
#    docker build . -t {{image_name}}:{{tag}} {{flags}}
#    rm -rf extproc

up:
    just containerize compose
    docker compose up

call path="/" method="GET" *flags="": 
    curl -X {{method}} localhost:8080{{path}} {{flags}} -vvv -s | jq .

down:
    docker compose down
