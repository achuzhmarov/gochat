package dal

import (
	"gochat/model"
	"gochat/test"

	"fmt"
	"testing"

	"gopkg.in/mgo.v2/bson"
)

var testMessageDao = MessageDao{"test_messages"}

var testUser *model.User
var testFriend1 *model.User
var testFriend2 *model.User

var testMessageTalkId1 bson.ObjectId
var testMessageTalkId2 bson.ObjectId

func initUsers() {
	testUserDao.ClearCollection()
	testUser = test.GenerateBaseUser()
	testFriend1 = test.GenerateBaseUserFriend1()
	testFriend2 = test.GenerateBaseUserFriend2()
	testUserDao.Upsert(testUser)
	testUserDao.Upsert(testFriend1)
	testUserDao.Upsert(testFriend2)

	testMessageTalkId1 = testFriend1.Talks[0].Id
	testMessageTalkId2 = testFriend2.Talks[0].Id
}

func TestUpsert_CreateMessage_NoErrors(t *testing.T) {
	initUsers()
	testMessageDao.ClearCollection()
	defer testMessageDao.ClearCollection()

	message := test.GenerateBaseMessage(testMessageTalkId1)

	_, err := testMessageDao.Upsert(message)
	test.TestForError(t, err)

	fmt.Println(message)

	messages, err := testMessageDao.GetAll()
	test.TestForError(t, err)

	fmt.Println(messages)
}

func TestUpsert_UpdateMessage_NoErrors(t *testing.T) {
	initUsers()
	testMessageDao.ClearCollection()
	defer testMessageDao.ClearCollection()

	message := test.GenerateBaseMessage(testMessageTalkId1)

	_, err := testMessageDao.Upsert(message)
	test.TestForError(t, err)

	fmt.Println(message)

	word3 := model.CreateNewWord("word3", "!", "W")
	message.Words = append(message.Words, word3)

	_, err = testMessageDao.Upsert(message)
	test.TestForError(t, err)

	messages, err := testMessageDao.GetAll()
	test.TestForError(t, err)

	fmt.Println(messages)
}

func TestDelete_DeleteMessage_NoErrors(t *testing.T) {
	initUsers()
	testMessageDao.ClearCollection()
	defer testMessageDao.ClearCollection()

	message := test.GenerateBaseMessage(testMessageTalkId1)

	_, err := testMessageDao.Upsert(message)
	test.TestForError(t, err)

	fmt.Println(message)

	err = testMessageDao.Delete(message)
	test.TestForError(t, err)

	messages, err := testMessageDao.GetAll()
	test.TestForError(t, err)

	fmt.Println(messages)
}

func TestGetMessagesByTalk_GetMessages_NoErrors(t *testing.T) {
	initUsers()
	testMessageDao.ClearCollection()
	defer testMessageDao.ClearCollection()

	message11 := test.GenerateBaseMessage(testMessageTalkId1)
	message12 := test.GenerateBaseMessage(testMessageTalkId1)
	message12.Author = "user_friend1"
	message21 := test.GenerateBaseMessage(testMessageTalkId2)
	message22 := test.GenerateBaseMessage(testMessageTalkId2)
	message22.Author = "user_friend2"

	_, err := testMessageDao.Upsert(message11)
	test.TestForError(t, err)

	_, err = testMessageDao.Upsert(message12)
	test.TestForError(t, err)

	_, err = testMessageDao.Upsert(message21)
	test.TestForError(t, err)

	_, err = testMessageDao.Upsert(message22)
	test.TestForError(t, err)

	fmt.Println(message11)
	fmt.Println(message12)
	fmt.Println(message21)
	fmt.Println(message22)

	messages1, err := testMessageDao.GetAllByTalk(testMessageTalkId1)
	test.TestForError(t, err)

	messages2, err := testMessageDao.GetAllByTalk(testMessageTalkId2)
	test.TestForError(t, err)

	fmt.Println(messages1)
	fmt.Println(messages2)
}

func TestMessageDao_ClearTalks(t *testing.T) {
	testUserDao.ClearCollection()
}
