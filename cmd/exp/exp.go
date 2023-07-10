package main

import (
	"fmt"

	"github.com/AkifhanIlgaz/vocab-builder/parser"
)

func main() {
	wordInfo, _ := parser.ParseWord("https://www.oxfordlearnersdictionaries.com/definition/english/wall_1")
	fmt.Printf("%+v\n", wordInfo)
}
