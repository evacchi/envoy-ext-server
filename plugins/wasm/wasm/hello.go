package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	headers := map[string]string{}
	for _, arg := range os.Args[1:] {
		kv := strings.Split(arg, "=")
		headers[kv[0]] = kv[1]
	}
	onRequest(headers, os.Stdin)
}

func onRequest(reqHeaders map[string]string, reqBody io.Reader) {
	if h, ok := reqHeaders["x-wasm"]; ok {
		if h == "append" {
			bytes, err := io.ReadAll(reqBody)
			if err != nil {
				log.Fatal(err)
			}
			kv := map[string]string{}
			err = json.Unmarshal(bytes, &kv)
			if err != nil {
				log.Fatal(err)
			}
			kv["trailer"] = "...and that's all folks"

			bytes, err = json.Marshal(kv)
			if err != nil {
				log.Fatal(err)
			}
			os.Stdout.Write(bytes)
		}
	}
	return
}
