package parser

type PartOfSpeech int

const (
	Noun PartOfSpeech = iota
	Pronoun
	Verb
	Adjective
	Adverb
	Preposition
	Conjunction
	Interjection
	InvalidPartOfSpeech
)

type CEFRLevel int

const (
	A1 PartOfSpeech = iota
	A2
	B1
	B2
	C1
	C2
)
