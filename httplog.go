package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	addr     string
	location string
	cors     string
	known    = make(map[string]int)
)

func main() {
	flag.StringVar(&addr, "addr", ":8080", "address on which to listen")
	flag.StringVar(&location, "Hlocation", "", "Location field of the header")
	flag.StringVar(&cors, "HCrossOrigin", "*", "Location field of the header")
	flag.Parse()

	http.HandleFunc("/", logAll)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func logAll(rw http.ResponseWriter, req *http.Request) {
	summary := "<empty>"
	if req.ContentLength > 0 {
		data, err := ioutil.ReadAll(req.Body)
		if err != nil {
			panic(err)
		}

		dataStr := string(data)
		maxIdx := min(len(dataStr)-1, 140)
		summary = dataStr[:maxIdx]
	}

	entry := fmt.Sprintf("method=%s path=%s, body=%s, raddr=%s", req.Method, req.URL.Path, summary, req.RemoteAddr)

	if c, ok := known[entry]; !ok {
		log.Print(entry)
		known[entry] = 0
	} else {
		known[entry] = c + 1
		log.Printf("%d times : `%s...`", c+1, entry[:min(len(entry)-1, 30)])
	}

	if location != "" {
		rw.Header().Set("Location", location)
	}
	if location != "" {
		rw.Header().Set("Access-Control-Allow-Origin", cors)
	}
	rw.WriteHeader(http.StatusOK)
}

func min(n, m int) int {
	if n < m {
		return n
	}
	return m
}
func max(n, m int) int {
	if n > m {
		return n
	}
	return m
}