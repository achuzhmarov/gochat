package model

import (
	"fmt"

	"gopkg.in/mgo.v2/bson"
)

type Talk struct {
	Id          bson.ObjectId `json:"id" bson:"_id"`
	Users       []string      `json:"users"`
	Size        int           `json:"size"`
	LastMessage *Message      `json:"last_message" bson:"last_message"`
	HasUnread   bool          `json:"has_unread" bson:"has_unread"`
	CipheredNum int           `json:"ciphered_num" bson:"ciphered_num"`
}

type Talks []*Talk

func (t *Talk) GetId() bson.ObjectId {
	return t.Id
}

func (t *Talk) SetId(id bson.ObjectId) {
	t.Id = id
}

func CreateTalk(id bson.ObjectId, users []string) *Talk {
	return &Talk{
		id,
		users,
		len(users),
		nil,
		false,
		0,
	}
}

func (t Talks) String() string {
	result := ""

	for _, talk := range t {
		result = result + fmt.Sprint(talk)
	}

	return result
}
