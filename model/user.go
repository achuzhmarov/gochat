package model

import (
	"fmt"

	"gopkg.in/mgo.v2/bson"
)

type UserAuth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type User struct {
	Id          bson.ObjectId `json:"-" bson:"_id"`
	Login       string        `json:"login"`
	Password    string        `json:"-"`
	Email       string        `json:"-"`
	Friends     []string      `json:"friends"`
	Suggestions int           `json:"suggestions"`
	Rating      int           `json:"rating"`
	Talks       Talks         `json:"talks"`
}

type Users []*User

func (u *User) GetId() bson.ObjectId {
	return u.Id
}

func (u *User) SetId(id bson.ObjectId) {
	u.Id = id
}

func CreateUser(
	login string,
	password string,
	email string,
) *User {
	return &User{
		"",
		login,
		password,
		email,
		[]string{},
		10,
		0,
		Talks{},
	}
}

func (u Users) String() string {
	result := ""

	for _, user := range u {
		result = result + fmt.Sprint(user)
	}

	return result
}
