package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {

	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte(fmt.Sprintf("Hello from %s 👋👋👋\n", hostname)))
	})
	http.ListenAndServe(":8080", nil)
}
