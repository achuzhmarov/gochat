package logic

import (
	"gochat/dal"
	"gochat/model/wordType"
	"gochat/test"

	"testing"
	"fmt"
)

var testMessageDao = dal.MessageDao{"test_messages"}

var TestMessageService = messageService{testMessageDao, testUserDao}

func TestCreateMessage_okMessage_noErrors(t *testing.T) {
	initUsers()

	testMessageDao.ClearCollection()
	defer testMessageDao.ClearCollection()

	message := test.GenerateBaseMessage(testMessageTalk1.Id)
	_, err := TestMessageService.CreateMessage(testFriend1.Login, message)
	test.TestForError(t, err)
	fmt.Println(message)

	messages, err := testMessageDao.GetAll()
	test.TestForError(t, err)
	fmt.Println(messages)

	users, err := testUserDao.GetAll()
	test.TestForError(t, err)
	fmt.Println(users)
}

func TestUpdateMessage_okMessage_noErrors(t *testing.T) {
	initUsers()

	testMessageDao.ClearCollection()
	defer testMessageDao.ClearCollection()

	message := test.GenerateBaseMessage(testMessageTalk1.Id)
	_, err := TestMessageService.CreateMessage(testFriend1.Login, message)
	test.TestForError(t, err)
	testMessageTalk1.LastMessage = message
	fmt.Println(message)

	message.Words[0].WordType = wordType.Success
	message.Words[1].WordType = wordType.Success
	_, err = TestMessageService.DecipherMessage(testFriend1.Login, message)
	test.TestForError(t, err)
	fmt.Println(message)

	messages, err := testMessageDao.GetAll()
	test.TestForError(t, err)
	fmt.Println(messages)

	users, err := testUserDao.GetAll()
	test.TestForError(t, err)
	fmt.Println(users)
}

func TestGetMessagesByTalk_ok_markedAsReaded(t *testing.T) {
	initUsers()

	testMessageDao.ClearCollection()
	defer testMessageDao.ClearCollection()

	message := test.GenerateBaseMessage(testMessageTalk1.Id)
	_, err := TestMessageService.CreateMessage(testFriend1.Login, message)
	test.TestForError(t, err)
	fmt.Println(message)

	messages, err := TestMessageService.GetMessagesByTalk(testFriend1.Login, testMessageTalk1.Id.Hex())
	test.TestForError(t, err)
	fmt.Println(messages)

	user, err := testUserDao.GetByLogin(testFriend1.Login)
	test.TestForError(t, err)
	if user.Talks[0].HasUnread {
		t.Error("Has unread talks")
	}
	fmt.Println(user)
}

func TestMessageService_ClearUsers(t *testing.T) {
	testUserDao.ClearCollection()
}
