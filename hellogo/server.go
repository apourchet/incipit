package main

import (
	"fmt"
	"net/http"
)

func helloGo(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(rw, "Hello Go!")
}

func main() {
	http.HandleFunc("/hellogo", helloGo)
	http.ListenAndServe(":8080", nil)
}
