package main

import (
	_ "embed"
	"log"
	"net/http"

	"github.com/flowchartsman/swaggerui"
)

//go:embed api/*.yaml

var spec []byte

func main() {
	log.SetFlags(0)
	http.Handle("/swagger/", http.StripPrefix("/swagger", swaggerui.Handler(spec)))
	log.Println("serving on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
