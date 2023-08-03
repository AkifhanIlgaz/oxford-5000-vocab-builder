package parser

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/AkifhanIlgaz/vocab-builder/models"
	"github.com/PuerkitoBio/goquery"
)

func ParseWord(wordUrl string) (models.WordInfo, error) {
	var wordInfo models.WordInfo

	resp, err := http.Get(wordUrl)
	if err != nil {
		return wordInfo, fmt.Errorf("parsing word: %w", err)
	}
	defer resp.Body.Close()

	document, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return wordInfo, fmt.Errorf("parsing word: %w", err)
	}

	mainContainer := document.Find("#entryContent")

	parseHeader(mainContainer.Find(".webtop"), &wordInfo)
	parseDefinitions(mainContainer.Find("ol.senses_multiple").First().Find("li.sense"), &wordInfo)
	parseDefinitions(mainContainer.Find("ol.sense_single").First().Find("li.sense"), &wordInfo)
	parseIdioms(mainContainer.Find("div.idioms .idm-g"), &wordInfo)

	return wordInfo, nil
}

func parseHeader(mainContainer *goquery.Selection, wordInfo *models.WordInfo) {
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

func parseDefinitions(mainContainer *goquery.Selection, wordInfo *models.WordInfo) {

	mainContainer.Each(func(i int, s *goquery.Selection) {
		var definition models.Definition

		s.Find("span.def").Each(func(i int, s *goquery.Selection) {
			definition.Meaning = s.Text()
		})

		s.Find("ul.examples > li span.x").Each(func(i int, s *goquery.Selection) {
			html, _ := s.Html()
			definition.Examples = append(definition.Examples, html)
		})
		fmt.Println(definition)
		wordInfo.Definitions = append(wordInfo.Definitions, definition)

	})

	fmt.Println(wordInfo.Definitions)
}

func parseIdioms(mainContainer *goquery.Selection, wordInfo *models.WordInfo) {
	mainContainer.Each(func(i int, s *goquery.Selection) {
		var idiom models.Idiom
		idiom.Usage = s.Find("div.top-container").Text()

		s.Find(`ol[class^="sense"] li.sense`).Each(func(i int, s *goquery.Selection) {
			var definition models.Definition
			definition.Meaning = s.Find("span.def").Text()

			s.Find("ul.examples li span.x").Each(func(i int, s *goquery.Selection) {
				definition.Examples = append(definition.Examples, s.Text())
			})

			idiom.Definitions = append(idiom.Definitions, definition)
		})
		wordInfo.Idioms = append(wordInfo.Idioms, idiom)
		idiom = models.Idiom{}
	})
}
