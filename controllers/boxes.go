package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/AkifhanIlgaz/vocab-builder/context"
	"github.com/AkifhanIlgaz/vocab-builder/models"
	"github.com/go-chi/chi/v5"
)

type BoxController struct {
	BoxService  *models.BoxService
	WordService *models.WordService
	// TODO: Add other services if necessary
}

// TODO: Delete this function
func (bc *BoxController) GetWordBox(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())

	words, _ := bc.BoxService.GetWordBox(user.Id)

	bc.encode(w, words[:3])
}

// OK
func (bc *BoxController) NewWordBox(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())

	err := bc.BoxService.CreateWordBox(user.Id)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	// TODO: Redirect to todays words or home page
	fmt.Fprint(w, "redirected to box/today")
}

// OK
func (bc *BoxController) GetTodaysWords(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	words, err := bc.BoxService.GetTodaysWords(user.Id)

	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	var wordInfos []*models.WordInfo

	for i := 0; i < 3; i++ {
		wordInfo, _ := bc.WordService.GetWord(words[i].Id)
		wordInfos = append(wordInfos, wordInfo)
	}

	bc.encode(w, wordInfos)
}

// OK
func (bc *BoxController) LevelUp(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	// Error handling
	wordId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid word id", http.StatusBadRequest)
		return
	}

	if err := bc.BoxService.LevelUp(user.Id, wordId); err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "llevel up !")
}

// OK
func (bc *BoxController) LevelDown(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	// Error handling
	wordId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid word id", http.StatusBadRequest)
		return
	}

	if err := bc.BoxService.LevelDown(user.Id, wordId); err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "llevel down !")
}

func (bc *BoxController) encode(w http.ResponseWriter, data any) {
	enc := json.NewEncoder(w)
	// TODO: Should I set indent ?
	enc.SetIndent("", "  ")
	enc.Encode(data)
}
