package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AkifhanIlgaz/vocab-builder/parser"
	"github.com/go-chi/chi/v5"
)

const baseUrl = "https://www.oxfordlearnersdictionaries.com/definition/english/"

func main() {
	router := chi.NewRouter()

	router.Get("/{word}", func(w http.ResponseWriter, r *http.Request) {
		wordInfo, err := parser.ParseWord(baseUrl + chi.URLParam(r, "word"))
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}

		x, _ := json.MarshalIndent(wordInfo, "", "\t")

		w.Write(x)
	})

	fmt.Println("Serving")
	http.ListenAndServe(":3000", router)
}
