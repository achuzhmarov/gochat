package model

import (
	"gochat/model/wordType"

	"fmt"
)

type Word struct {
	Text         string            `json:"text" bson:"text"`
	Additional   string            `json:"additional"`
	CipheredText string            `json:"ciphered_text" bson:"ciphered_text"`
	WordType     wordType.WordType `json:"word_type" bson:"word_type"`
}

type Words []*Word

func CreateNewWord(text string, additional string, cipheredText string) *Word {
	return &Word{
		text,
		additional,
		cipheredText,
		wordType.New,
	}
}

func CreateWord(text string, additional string, cipheredText string, wordType wordType.WordType) *Word {
	return &Word{
		text,
		additional,
		cipheredText,
		wordType,
	}
}

func (w Words) String() string {
	result := ""

	for _, word := range w {
		result = result + fmt.Sprint(word)
	}

	return result
}
