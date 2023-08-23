package main

import (
	"fmt"
	"net/http"
)

func main() {
	var a int
	http.HandleFunc("/increment", func(w http.ResponseWriter, r *http.Request) {
		a++
	})
	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintln(a)))
	})
	http.ListenAndServe(":2345", nil)
}
