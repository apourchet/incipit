package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/context"

	"github.com/apourchet/dummy/lib/etcd"
	"github.com/apourchet/dummy/lib/healthz"
	"github.com/apourchet/dummy/lib/utils"
)

func helloHermes(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(rw, "Hello Hermes!")
}

func main() {
	utils.Info("Starting the server on 8080")
	http.HandleFunc("/", helloHermes)

	etcd, err := etcd.GetK8sDefaultClient()
	if err != nil {
		log.Fatal(err)
	}

	utils.Info("Setting '/foo' key with 'bar' value")
	resp, err := etcd.Set(context.Background(), "/foo", "bar", nil)
	if err != nil {
		utils.Fatal(err)
	} else {
		utils.Info("Set is done. Metadata is %q\n", resp)
	}

	healthz.SpawnHealthCheck(healthz.DEFAULT_PORT)
	http.ListenAndServe(":8080", nil)
}
