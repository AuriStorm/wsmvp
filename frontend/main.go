package main

import (
	"fmt"
	"net/http"
)

func logHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf(
			"start request; method: %s; host: %s; url: %s; remote_addr: %s\n",
			r.Method, r.Host, r.URL, r.RemoteAddr,
		)
		h.ServeHTTP(w, r)
		fmt.Println("finish request;")
	}
}

func main() {
	http.Handle("/", logHandler(http.FileServer(http.Dir("./"))))
	fmt.Println("Start servser on localhost:3000")

	http.ListenAndServe(":3000", nil)
}
