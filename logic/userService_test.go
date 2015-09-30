package logic

import (
	"gochat/dal"
	"gochat/model"
	"gochat/test"

	"testing"
	"fmt"
)

var testUserDao = dal.UserDao{"test_users"}
var TestUserService = userService{testUserDao}

var testUser *model.User
var testFriend1 *model.User
var testFriend2 *model.User

var testMessageTalk1 *model.Talk
var testMessageTalk2 *model.Talk

//for real database
/*func testInitUsers(t *testing.T) {
	realUserDao := dal.UserDao{"users"}
	realUserDao.ClearCollection()
	realUserDao.Upsert(test.GenerateBaseUser())
	realUserDao.Upsert(test.GenerateBaseUserFriend1())
	realUserDao.Upsert(test.GenerateBaseUserFriend2())
}*/

func initUsers() {
	testUserDao.ClearCollection()
	testUser = test.GenerateBaseUser()
	testFriend1 = test.GenerateBaseUserFriend1()
	testFriend2 = test.GenerateBaseUserFriend2()
	testUserDao.Upsert(testUser)
	testUserDao.Upsert(testFriend1)
	testUserDao.Upsert(testFriend2)

	testUserDao.Upsert(test.GenerateEmptyUser("user333"))
	testUserDao.Upsert(test.GenerateEmptyUser("user444"))
	testUserDao.Upsert(test.GenerateEmptyUser("user345"))
	testUserDao.Upsert(test.GenerateEmptyUser("aaa"))

	testMessageTalk1 = testFriend1.Talks[0]
	testMessageTalk2 = testFriend2.Talks[0]
}

func TestGetNewFriends_ok_noErrors(t *testing.T) {
	initUsers()

	friends, err := TestUserService.GetNewFriends("user", "")
	test.TestForError(t, err)
	fmt.Println(friends)

	friends, err = TestUserService.GetNewFriends("user", "user")
	test.TestForError(t, err)
	fmt.Println(friends)

	friends, err = TestUserService.GetNewFriends("user", "3")
	test.TestForError(t, err)
	fmt.Println(friends)

	friends, err = TestUserService.GetNewFriends("user", "34")
	test.TestForError(t, err)
	fmt.Println(friends)
}

func TestMakeFriends_ok_noErrors(t *testing.T) {
	initUsers()

	_, err := TestUserService.MakeFriends("user333", "user444")
	test.TestForError(t, err)
	_, err = TestUserService.MakeFriends("user333", "user")
	test.TestForError(t, err)
	_, err = TestUserService.MakeFriends("user444", "user")
	test.TestForError(t, err)

	user333, err := TestUserService.GetUserByLogin("user333")
	user444, err := TestUserService.GetUserByLogin("user444")
	fmt.Println(user333.Friends)
	fmt.Println(user444.Friends)
	fmt.Println(user333.Talks)
	fmt.Println(user444.Talks)
}

func TestUserService_ClearUsers(t *testing.T) {
	testUserDao.ClearCollection()
}
