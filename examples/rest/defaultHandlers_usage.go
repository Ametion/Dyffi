package main

import (
	"net/http"

	"github.com/Ametion/dyffi"
)

type customHandler struct {
	message string
}

func (h *customHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(h.message))
}

func main() {
	engine := dyffi.NewDyffiEngine()
	
	engine.Get("/default", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("some  string"))
	})

	engine.Get("/default_handler", &customHandler{message: "some string"})


	engine.Run(":8080")
}

