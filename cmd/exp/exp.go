package main

import (
	"encoding/json"
	"fmt"

	"github.com/AkifhanIlgaz/vocab-builder/parser"
)

const (
	withIdioms    = "https://www.oxfordlearnersdictionaries.com/definition/english/about_2"
	withoutIdioms = "https://www.oxfordlearnersdictionaries.com/definition/english/across_2"
	diff          = "https://www.oxfordlearnersdictionaries.com/definition/english/reject_1"
)

func main() {
	wordInfo, _ := parser.ParseWord(withIdioms)

	x, _ := json.MarshalIndent(wordInfo, "", "\t")

	fmt.Println(string(x))
}
