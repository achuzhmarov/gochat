package test

import (
	"gochat/model"
	"gochat/model/cipher"

	"testing"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

func TestForError(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

var talk1Id = bson.NewObjectId()
var talk2Id = bson.NewObjectId()

func GenerateBaseUser() *model.User {
	user := model.CreateUser("user", getHash("user_pass"), "user@email")

	user.Friends = []string{"user1", "user2"}
	user.Rating = 15
	user.Suggestions = 10

	talk1 := generateBaseTalk(talk1Id, []string{"user", "user1"})
	talk2 := generateBaseTalk(talk2Id, []string{"user", "user2"})
	user.Talks = append(user.Talks, talk1, talk2)

	return user
}

func GenerateBaseUserFriend1() *model.User {
	user := model.CreateUser("user1", getHash("user1_pass"), "user1@email")

	user.Friends = []string{"user"}
	user.Rating = 15
	user.Suggestions = 10

	talk := generateBaseTalk(talk1Id, []string{"user", "user1"})
	user.Talks = append(user.Talks, talk)

	return user
}

func GenerateBaseUserFriend2() *model.User {
	user := model.CreateUser("user2", getHash("user2_pass"), "user2@email")

	user.Friends = []string{"user"}
	user.Rating = 15
	user.Suggestions = 10

	talk := generateBaseTalk(talk2Id, []string{"user", "user2"})
	user.Talks = append(user.Talks, talk)

	return user
}

func GenerateEmptyUser(userName string) *model.User {
	return model.CreateUser(userName, getHash("pass"), "user@email")
}

func GenerateBaseMessage(talkId bson.ObjectId) *model.Message {
	word1 := model.CreateNewWord("word1", "?", "What")
	word2 := model.CreateNewWord("word2", "!!", "Wow")
	words := model.Words{word1, word2}
	return model.CreateMessage(talkId, "user", false, cipher.FirstLetterCipher, words)
}

func generateBaseTalk(talkId bson.ObjectId, users []string) *model.Talk {
	return model.CreateTalk(talkId, users)
}

func getHash(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(hashedPassword)
}
