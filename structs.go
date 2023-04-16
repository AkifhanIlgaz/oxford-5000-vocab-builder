package main

type WordInfo struct {
	Word         string    `json:"Word"`
	PartOfSpeech string    `json:"PartOfSpeech"`
	CEFRLevel    string    `json:"CEFRLevel"`
	Meanings     []Meaning `json:"Meanings"`
	Idioms       []Idiom   `json:"Idioms"`
	Url          string    `json:"URL"`
}

type Idiom struct {
	Idm        string   `json:"idm"`
	Definition string   `json:"def"`
	Examples   []string `json:"examples"`
}

type Meaning struct {
	CF         string   `json:"cf"`
	Definition string   `json:"Definition"`
	Examples   []string `json:"Examples"`
}
