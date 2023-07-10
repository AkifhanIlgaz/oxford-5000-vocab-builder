package main

import (
	"encoding/json"
	"fmt"

	"github.com/AkifhanIlgaz/vocab-builder/parser"
)

func main() {
	wordInfo, _ := parser.ParseWord("https://www.oxfordlearnersdictionaries.com/definition/english/abandon_1")

	x, _ := json.MarshalIndent(wordInfo, "", "\t")

	fmt.Println(string(x))
}
