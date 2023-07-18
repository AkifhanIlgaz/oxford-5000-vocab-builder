package controllers

import (
	"net/http"

	"github.com/AkifhanIlgaz/vocab-builder/models"
)

type BoxController struct {
	BoxService *models.BoxService
	// TODO: Add other services if necessary
}

func (bc *BoxController) NewWordBox(w http.ResponseWriter, r *http.Request) {

}

func (bc *BoxController) GetTodaysWords(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get todays words"))
}

func (bc *BoxController) GetWordByLevel(w http.ResponseWriter, r *http.Request) {
	// TODO: Extract user from request's context
}

// TODO: Change function names. There is duplicate functions on different modules
func (bc *BoxController) LevelUp(w http.ResponseWriter, r *http.Request) {

}

func (bc *BoxController) LevelDown(w http.ResponseWriter, r *http.Request) {

}
