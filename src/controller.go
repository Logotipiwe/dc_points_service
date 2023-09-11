package main

import (
	"net/http"
	"os"
)

func main() {
	err := InitDb()
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		println("/auth")
	})

	println("Ready")
	err = http.ListenAndServe(":"+os.Getenv("CONTAINER_PORT"), nil)
	println("Server up!")
	if err != nil {
		panic("Lol server fell")
	}
}
