package controllers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/AkifhanIlgaz/vocab-builder/context"
	"github.com/AkifhanIlgaz/vocab-builder/models"
	"github.com/go-chi/chi/v5"
)

type BoxController struct {
	BoxService  *models.BoxService
	WordService *models.WordService
}

// OK
func (bc *BoxController) GetTodaysWords(w http.ResponseWriter, r *http.Request) {
	uid := context.Uid(r.Context())
	words, err := bc.BoxService.GetTodaysWords(uid)
	if len(words) == 0 {
		if err := bc.BoxService.CreateWordBox(uid); err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		words, err = bc.BoxService.GetTodaysWords(uid)
	}

	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	var wordInfos []*models.WordInfo

	for _, i := range rand.Perm(len(words))[:10] {
		wordInfo, _ := bc.WordService.GetWord(words[i].Id)
		wordInfos = append(wordInfos, wordInfo)
	}

	bc.encode(w, wordInfos)
}

// OK
func (bc *BoxController) LevelUp(w http.ResponseWriter, r *http.Request) {
	uid := context.Uid(r.Context())
	wordId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid word id", http.StatusBadRequest)
		return
	}

	if err := bc.BoxService.LevelUp(uid, wordId); err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "llevel up !")
}

func (bc *BoxController) LevelDown(w http.ResponseWriter, r *http.Request) {
	uid := context.Uid(r.Context())
	// Error handling
	wordId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid word id", http.StatusBadRequest)
		return
	}

	if err := bc.BoxService.LevelDown(uid, wordId); err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "llevel down !")
}

func (bc *BoxController) encode(w http.ResponseWriter, data any) {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(data)
}
