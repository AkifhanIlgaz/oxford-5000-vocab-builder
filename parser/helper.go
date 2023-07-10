package parser

// TODO: Rename this function
func convertPosConstant(pos string) PartOfSpeech {
	switch pos {
	case "noun":
		return Noun
	case "pronoun":
		return Pronoun
	case "verb":
		return Verb
	case "adjective":
		return Adjective
	case "adverb":
		return Adverb
	case "preposition":
		return Preposition
	case "conjunction":
		return Conjunction
	case "interjection":
		return Interjection
	default:
		return InvalidPartOfSpeech
	}
}
