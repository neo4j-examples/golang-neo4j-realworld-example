package main

import "net/http"

func main() {

	server := http.NewServeMux()
	server.HandleFunc("/users", func(writer http.ResponseWriter, request *http.Request) {

	})

	if err := http.ListenAndServe(":3000", server); err != nil {
		panic(err)
	}

}
