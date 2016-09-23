package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/context"

	"github.com/apourchet/dummy/lib/etcd"
)

func helloGo(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(rw, "Hello Go!")
}

func main() {
	log.Println("Starting the server on 8080")
	http.HandleFunc("/hellogo", helloGo)

	etcd, err := etcd.GetK8sDefaultClient()
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Setting '/foo' key with 'bar' value")
	resp, err := etcd.Set(context.Background(), "/foo", "bar", nil)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Set is done. Metadata is %q\n", resp)
	}

	http.ListenAndServe(":8080", nil)
}
