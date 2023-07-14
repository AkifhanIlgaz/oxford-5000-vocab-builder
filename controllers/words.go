package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/AkifhanIlgaz/vocab-builder/models"
	"github.com/go-chi/chi/v5"
)

type WordsController struct {
	WordService *models.WordService
}

func (wc WordsController) WordWithId(w http.ResponseWriter, r *http.Request) {
	wordId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid word ID", http.StatusNotFound)
		return
	}

	word, err := wc.WordService.GetWord(wordId)
	if err != nil {
		http.Error(w, "cannot found word", http.StatusNotFound)
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(word)
}
