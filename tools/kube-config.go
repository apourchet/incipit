package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
)

func usage() {
	fmt.Println("Must have at least one configuration file and one k8s-template as arguments")
}

func main() {
	if len(os.Args) < 3 {
		usage()
		return
	}

	config := make(map[string]string)
	configBytes, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read config file %s\n", os.Args[1])
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse config file %s\n", os.Args[1])
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
		return
	}
	for i := 2; i < len(os.Args); i++ {
		fname := os.Args[i]
		templateBytes, err := ioutil.ReadFile(fname)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read template file: %s\n", fname)
			fmt.Fprintf(os.Stderr, err.Error())
			continue
		}
		tmpl, err := template.New(fname).Parse(string(templateBytes))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse template file: %s\n", fname)
			fmt.Fprintf(os.Stderr, err.Error())
			continue
		}
		err = tmpl.Execute(os.Stdout, config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to Execute template %s\n", fname)
			fmt.Fprintf(os.Stderr, err.Error())
			continue
		}
		fmt.Println("---")
	}
}