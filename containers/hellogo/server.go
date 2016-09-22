package main

import (
	"fmt"
	"net/http"
)

func helloGo(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(rw, "Hello Go!")
}

func main() {
	fmt.Println("Starting the server on 8080")
	http.HandleFunc("/hellogo", helloGo)
	http.ListenAndServe(":8080", nil)
}
