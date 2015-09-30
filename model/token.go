package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Token struct {
	Id       bson.ObjectId `json:"id" bson:"_id"`
	Value    string
	ExpireAt time.Time `bson:"expire_at"`
}

type Tokens []*Token

func (t *Token) GetId() bson.ObjectId {
	return t.Id
}

func (t *Token) SetId(id bson.ObjectId) {
	t.Id = id
}

func CreateToken(value string, expireAt time.Time) *Token {
	return &Token{"", value, expireAt}
}
