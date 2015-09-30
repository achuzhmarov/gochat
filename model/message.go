package model

import (
	"gochat/model/cipher"
	"gochat/model/wordType"

	"fmt"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Message struct {
	Id         bson.ObjectId     `json:"id" bson:"_id"`
	TalkId     bson.ObjectId     `json:"talk_id" bson:"_talk_id"`
	Author     string            `json:"author"`
	Words      Words             `json:"words"`
	Deciphered bool              `json:"deciphered"`
	CipherType cipher.CipherType `json:"cipher_type" bson:"cipher_type"`
	Timestamp  time.Time         `json:"timestamp"`
}

type Messages []*Message

func (m *Message) GetId() bson.ObjectId {
	return m.Id
}

func (m *Message) SetId(id bson.ObjectId) {
	m.Id = id
}

func CreateMessage(talkId bson.ObjectId, author string, deciphered bool,
	cipherType cipher.CipherType, words Words) *Message {

	return &Message{
		"",
		talkId,
		author,
		words,
		deciphered,
		cipherType,
		time.Now(),
	}
}

func (m *Message) CheckDeciphered() bool {
	for _, word := range m.Words {
		if word.WordType == wordType.New {
			return false
		}
	}

	m.Deciphered = true
	return true
}

func (m Messages) String() string {
	result := ""

	for _, message := range m {
		result = result + fmt.Sprint(message)
	}

	return result
}
