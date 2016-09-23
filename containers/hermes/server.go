package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/context"

	"github.com/coreos/etcd/client"
)

func helloGo(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(rw, "Hello Go!")
}

func main() {
	fmt.Println("Starting the server on 8080")
	http.HandleFunc("/hellogo", helloGo)

	etcdClientAddr := "http://" + os.Getenv("ETCD_CLIENT_SERVICE_HOST") + ":" + os.Getenv("ETCD_CLIENT_SERVICE_PORT")
	cfg := client.Config{
		Endpoints:               []string{etcdClientAddr},
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}

	c, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	kapi := client.NewKeysAPI(c)
	log.Print("Setting '/foo' key with 'bar' value")
	resp, err := kapi.Set(context.Background(), "/foo", "bar", nil)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Set is done. Metadata is %q\n", resp)
	}

	http.ListenAndServe(":8080", nil)
}
