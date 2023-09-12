package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	err := InitDb()
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		println("/auth")
		fmt.Fprintf(w, "Hi!")
	})

	println("Ready")
	err = http.ListenAndServe(":"+os.Getenv("CONTAINER_PORT"), nil)
	println("Server up!")
	if err != nil {
		panic("Lol server fell")
	}
}
