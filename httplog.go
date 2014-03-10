package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var (
	addr    string
	maxLen  int
	headers = make(HeaderValue)
	body    string
	known   = make(map[string]int)
)

func main() {
	flag.StringVar(&addr, "addr", ":8080", "address on which to listen")

	flag.Var(headers, "H", "comma separated list of a header key, followed by values. Invoke many times")
	flag.IntVar(&maxLen, "MaxLen", 140, "Max length of printed strings")
	flag.StringVar(&body, "Body", "", "Body to put in response")

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

	for headKey, headVal := range headers {
		rw.Header()[headKey] = headVal
	}

	rw.WriteHeader(http.StatusOK)
	fmt.Fprint(rw, body)
}

type HeaderValue map[string][]string

func (h HeaderValue) String() string {
	return "\"key, val1, val2, ..., valn\""
}

func (h HeaderValue) Set(s string) error {
	vals := strings.Split(s, ",")
	key := vals[0]
	vals = vals[1:]
	for i, each := range vals {
		vals[i] = strings.TrimSpace(each)
	}
	h[key] = vals
	return nil
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
