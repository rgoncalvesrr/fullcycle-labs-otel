package main

import "net/http"

func main() {
	mux := http.NewServeMux()

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		return
	}
}
