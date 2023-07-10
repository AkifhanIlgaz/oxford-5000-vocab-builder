package parser

type WordInfo struct {
	Word        string       `json:"word"`
	Header      Header       `json:"header"`
	Definitions []Definition `json:"definitions"`
	Idioms      []Idiom      `json:"idioms"`
}

type Header struct {
	Audio        []string     `json:"audio"`
	PartOfSpeech PartOfSpeech `json:"partOfSpeech"`
	CEFRLevel    CEFRLevel    `json:"CEFRLevel"`
}

type Definition struct {
	Meaning  string   `json:"meaning"`
	Examples []string `json:"examples"`
}

type Idiom struct {
	Usage      string     `json:"usage"`
	Definition Definition `json:"definition"`
}
