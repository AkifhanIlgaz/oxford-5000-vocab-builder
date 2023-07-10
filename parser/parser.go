package parser

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// TODO: Create custom client for HTTP requests
// TODO: Create custom errors

func ParseWord(wordUrl string) (WordInfo, error) {
	// TODO: Load HTML and select main container
	// Pass container to other functions
	var wordInfo WordInfo

	resp, err := http.Get(wordUrl)
	if err != nil {
		return wordInfo, fmt.Errorf("parsing word: %w", err)
	}
	defer resp.Body.Close()

	document, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return wordInfo, fmt.Errorf("parsing word: %w", err)
	}

	mainContainer := document.Find("#main-container")

	err = parseHeader(mainContainer, &wordInfo)
	if err != nil {
		return wordInfo, fmt.Errorf("parsing word: %w", err)
	}

	err = parseDefinitions(mainContainer, &wordInfo)
	if err != nil {
		return wordInfo, fmt.Errorf("parsing word: %w", err)
	}

	err = parseIdioms(mainContainer, &wordInfo)
	if err != nil {
		return wordInfo, fmt.Errorf("parsing word: %w", err)
	}

	return wordInfo, nil
}

func parseHeader(mainContainer *goquery.Selection, wordInfo *WordInfo) error {
	var wordHeader Header

	mainContainer.Find("span.pos").Each(func(i int, s *goquery.Selection) {
		wordHeader.PartOfSpeech = convertPosConstant(s.Text())
	})

	return nil
}

func parseDefinitions(mainContainer *goquery.Selection, wordInfo *WordInfo) error {
	return nil
}

func parseIdioms(mainContainer *goquery.Selection, wordInfo *WordInfo) error {
	return nil
}
