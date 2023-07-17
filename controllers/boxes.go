package controllers

import (
	"net/http"

	"github.com/AkifhanIlgaz/vocab-builder/models"
)

type BoxController struct {
	BoxService *models.BoxService
	// TODO: Add other services if necessary
}

func (bc *BoxController) GetTodaysWords(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get todays words"))

}

func (bc *BoxController) GetWordByLevel(w http.ResponseWriter, r *http.Request) {
	// TODO: Extract user from request's context
}
