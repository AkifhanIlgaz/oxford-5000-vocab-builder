package parser

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type WordInfo struct {
	Index  int    `json:"index"`
	Source string `json:"source"`
	Word   string `json:"word"`
	Header
	Definitions []Definition `json:"definitions"`
	Idioms      []Idiom      `json:"idioms"`
}

type Header struct {
	Audio struct {
		UK string `json:"UK" bson:"UK"`
		US string `json:"US" bson:"US"`
	} `json:"audio" bson:"audio"`
	PartOfSpeech string `json:"partOfSpeech" bson:"partOfSpeech"`
	CEFRLevel    string `json:"CEFRLevel" bson:"CEFRLevel"`
}

type Definition struct {
	Meaning  string   `json:"meaning"`
	Examples []string `json:"examples"`
}

type Idiom struct {
	Usage       string       `json:"usage"`
	Definitions []Definition `json:"definition"`
}

func ParseWord(wordUrl string) (WordInfo, error) {
	var wordInfo WordInfo

	req, err := http.NewRequest(http.MethodGet, wordUrl, nil)
	if err != nil {
		return wordInfo, fmt.Errorf("client: could not create request: %s\n", err)

	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return wordInfo, fmt.Errorf("client: error making http request: %s\n", err)
	}

	defer resp.Body.Close()

	document, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return wordInfo, fmt.Errorf("parsing word: %w", err)
	}

	if document == nil {
		return wordInfo, fmt.Errorf("parsing word: %w", err)
	}
	mainContainer := document.Find("#entryContent")

	parseHeader(mainContainer.Find(".webtop"), &wordInfo)
	parseDefinitions(mainContainer.Find("ol.senses_multiple").First().Find("li.sense"), &wordInfo)
	parseDefinitions(mainContainer.Find("ol.sense_single").First().Find("li.sense"), &wordInfo)
	parseIdioms(mainContainer.Find("div.idioms .idm-g"), &wordInfo)

	return wordInfo, nil
}

func parseHeader(mainContainer *goquery.Selection, wordInfo *WordInfo) {
	// HeadingWord
	mainContainer.Find(".headword").First().Each(func(i int, s *goquery.Selection) {
		wordInfo.Word = s.Text()
	})

	// Part Of Speech
	mainContainer.Find("span.pos").Each(func(i int, s *goquery.Selection) {
		wordInfo.Header.PartOfSpeech = s.Text()
	})

	// CEFR Level
	mainContainer.Find(".symbols span").First().Each(func(i int, s *goquery.Selection) {
		attr, _ := s.Attr("class")
		wordInfo.Header.CEFRLevel = strings.ToUpper(strings.Split(attr, "_")[1])
	})

	// Audio
	mainContainer.Find(`span.phonetics div > div`).Each(func(i int, s *goquery.Selection) {
		audioUrl, _ := s.Attr("data-src-mp3")

		// We don't need to check `pron-us` since there is only two possibilities
		if s.HasClass("pron-uk") {
			wordInfo.Header.Audio.UK = audioUrl
		} else {
			wordInfo.Header.Audio.US = audioUrl
		}

	})

}

func parseDefinitions(mainContainer *goquery.Selection, wordInfo *WordInfo) {
	mainContainer.Each(func(i int, s *goquery.Selection) {
		var definition Definition

		s.Find("span.def").Each(func(i int, s *goquery.Selection) {
			definition.Meaning = s.Text()
		})

		s.Find("ul.examples > li span.x").Each(func(i int, s *goquery.Selection) {
			html, _ := s.Html()
			definition.Examples = append(definition.Examples, html)
		})
		wordInfo.Definitions = append(wordInfo.Definitions, definition)

	})

}

func parseIdioms(mainContainer *goquery.Selection, wordInfo *WordInfo) {
	mainContainer.Each(func(i int, s *goquery.Selection) {
		var idiom Idiom
		idiom.Usage = s.Find("div.top-container").Text()

		s.Find(`ol[class^="sense"] li.sense`).Each(func(i int, s *goquery.Selection) {
			var definition Definition
			definition.Meaning = s.Find("span.def").Text()

			s.Find("ul.examples li span.x").Each(func(i int, s *goquery.Selection) {
				definition.Examples = append(definition.Examples, s.Text())
			})

			idiom.Definitions = append(idiom.Definitions, definition)
		})
		wordInfo.Idioms = append(wordInfo.Idioms, idiom)
		idiom = Idiom{}
	})
}
