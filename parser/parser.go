package parser

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// TODO: Create custom client for HTTP requests
// TODO: Create custom errors
// Should I return errors from parsing functions ?

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

	parseHeader(mainContainer.Find(".webtop"), &wordInfo)
	parseDefinitions(mainContainer, &wordInfo)
	parseIdioms(mainContainer, &wordInfo)

	return wordInfo, nil
}

func parseHeader(mainContainer *goquery.Selection, wordInfo *WordInfo) {
	// Part Of Speech
	mainContainer.Find("span.pos").Each(func(i int, s *goquery.Selection) {
		wordInfo.Header.PartOfSpeech = s.Text()
	})

	// CEFR Level
	mainContainer.Find(".symbols span").Each(func(i int, s *goquery.Selection) {
		attr, _ := s.Attr("class")
		if pos, ok := strings.CutPrefix(attr, "ox3ksym_"); ok {
			wordInfo.Header.CEFRLevel = strings.ToUpper(pos)
		}
	})

	// Audio
	mainContainer.Find(`span.phonetics div div`).Each(func(i int, s *goquery.Selection) {
		audioUrl, _ := s.Attr("data-src-mp3")

		// We don't need to check `pron-us` since there is only two possibilities
		if s.HasClass("pron-uk") {
			wordInfo.Header.Audio.UK = audioUrl
		} else {
			wordInfo.Header.Audio.US = audioUrl
		}
	})

}

func parseDefinitions(mainContainer *goquery.Selection, wordInfo *WordInfo) error {
	return nil
}

func parseIdioms(mainContainer *goquery.Selection, wordInfo *WordInfo) error {
	return nil
}
