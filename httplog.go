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
	maxLen   int
	body     string
	known    = make(map[string]int)
)

func main() {
	flag.StringVar(&addr, "addr", ":8080", "address on which to listen")
	flag.StringVar(&location, "HLocation", "", "Location field of the header")
	flag.StringVar(&cors, "HCrossOrigin", "*", "CrossOrigin field of the header")
	flag.IntVar(&maxLen, "MaxLen", 140, "Max length of printed strings")
	flag.StringVar(&body, "body", "", "Body to put in response")
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
		maxIdx := min(len(dataStr)-1, maxLen)
		summary = dataStr[:maxIdx]
	}

	entry := fmt.Sprintf("method=%s path=%s, body=%s, raddr=%s", req.Method, req.URL.Path, summary, req.RemoteAddr)

	if c, ok := known[entry]; !ok {
		log.Print(entry)
		known[entry] = 0
	} else {
		known[entry] = c + 1
		log.Printf("%d times : `%s...`", c+1, entry[:min(len(entry)-1, maxLen)])
	}

	if location != "" {
		rw.Header().Set("Location", location)
	}
	if cors != "" {
		rw.Header().Set("Access-Control-Allow-Origin", cors)
	}
	rw.WriteHeader(http.StatusOK)
	fmt.Fprint(rw, body)
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
