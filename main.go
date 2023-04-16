package main

import (
	"encoding/json"
	"fmt"
	"os"
)


const path string = "word_database/oxford5000words.json"

var words map[int]WordInfo

func main() {
	initDatabase()
	
	
}

// This function connects to the database. But for now, it just reads the json file and stores it in a map
// TODO: Convert this function to connect to a database
func initDatabase() {
	words = map[int]WordInfo{}

	f, _ := os.Open(path)
	defer f.Close()

	decoder := json.NewDecoder(f)

	decoder.Token()

	for decoder.More() {
		var word WordInfo
		decoder.Decode(&word)
		words[len(words)] = word
	}

	decoder.Token()
}
