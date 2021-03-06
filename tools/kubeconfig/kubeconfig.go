package main

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"io/ioutil"
	"os"
)

func usage() {
	fmt.Println("Must have at least one configuration file and one k8s-template as arguments")
}

func readFile(fname string) (string, error) {
	b, err := ioutil.ReadFile(fname)
	if err != nil {
		return "", err
	}
	return string(b), err
}

func base64Encode(s string) string {
	return string(b64.StdEncoding.EncodeToString([]byte(s)))
}

func loopOverInts(n float64) []int {
	arr := make([]int, int(n))
	for i := 0; i < int(n); i++ {
		arr[i] = i
	}
	return arr
}

func getFuncMap() template.FuncMap {
	return template.FuncMap{
		"include": readFile,
		"base64":  base64Encode,
		"loop":    loopOverInts,
	}
}

func objectsToString(m map[string]interface{}) error {
	for k, v := range m {
		cast, ok := v.(map[string]interface{})
		if ok {
			strBytes, err := json.Marshal(cast)
			if err != nil {
				return err
			}
			m[k] = string(strBytes)
		}
		cast1, ok := v.([]interface{})
		if ok {
			strBytes, err := json.Marshal(cast1)
			if err != nil {
				return err
			}
			m[k] = string(strBytes)
		}
	}
	return nil
}

func envOverride(m map[string]interface{}) {
	for k, _ := range m {
		if os.Getenv(k) != "" {
			m[k] = os.Getenv(k)
		}
	}
}

func main() {
	if len(os.Args) < 3 {
		usage()
		return
	}

	config := make(map[string]interface{})
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
	err = objectsToString(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse config file %s\n", os.Args[1])
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
		return
	}

	envOverride(config)

	for i := 2; i < len(os.Args); i++ {
		fname := os.Args[i]
		templateBytes, err := ioutil.ReadFile(fname)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read template file: %s\n", fname)
			fmt.Fprintf(os.Stderr, err.Error())
			continue
		}
		tmpl, err := template.New(fname).Funcs(getFuncMap()).Parse(string(templateBytes))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse template file: %s\n", fname)
			fmt.Fprintf(os.Stderr, err.Error())
			continue
		}
		var buf bytes.Buffer
		err = tmpl.Execute(&buf, config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to Execute template %s\n", fname)
			fmt.Fprintf(os.Stderr, err.Error())
			continue
		}
		str := html.UnescapeString(buf.String())
		fmt.Fprintf(os.Stdout, str)
		// fmt.Println("---")
	}
}
