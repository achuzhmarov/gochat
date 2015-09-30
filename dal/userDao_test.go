package dal

import (
	"gochat/test"

	"fmt"
	"testing"
)

var testUserDao = UserDao{"test_users"}

func TestUpsert_CreateUser_NoErrors(t *testing.T) {
	testUserDao.ClearCollection()
	defer testUserDao.ClearCollection()

	user := test.GenerateBaseUser()

	_, err := testUserDao.Upsert(user)
	test.TestForError(t, err)

	fmt.Println(user)

	users, err := testUserDao.GetAll()
	test.TestForError(t, err)

	fmt.Println(users)
}

func TestUpsert_UpdateUser_NoErrors(t *testing.T) {
	testUserDao.ClearCollection()
	defer testUserDao.ClearCollection()

	user := test.GenerateBaseUser()

	_, err := testUserDao.Upsert(user)
	test.TestForError(t, err)

	fmt.Println(user)

	user.Suggestions = 200
	user.Rating = 300
	user.Friends = append(user.Friends, "user_friend26")
	user.Email = "user@emailUpdated"

	_, err = testUserDao.Upsert(user)
	test.TestForError(t, err)

	users, err := testUserDao.GetAll()
	test.TestForError(t, err)

	fmt.Println(users)
}

func TestDelete_DeleteUser_NoErrors(t *testing.T) {
	testUserDao.ClearCollection()
	defer testUserDao.ClearCollection()

	user := test.GenerateBaseUser()

	_, err := testUserDao.Upsert(user)
	test.TestForError(t, err)

	fmt.Println(user)

	err = testUserDao.Delete(user)
	test.TestForError(t, err)

	users, err := testUserDao.GetAll()
	test.TestForError(t, err)

	fmt.Println(users)
}

func TestGetUserByLogin_GetUser_NoErrors(t *testing.T) {
	testUserDao.ClearCollection()
	defer testUserDao.ClearCollection()

	user := test.GenerateBaseUser()
	user2 := test.GenerateBaseUser()
	user2.Login = "test_login"

	_, err := testUserDao.Upsert(user)
	test.TestForError(t, err)

	_, err = testUserDao.Upsert(user2)
	test.TestForError(t, err)

	fmt.Println(user)
	fmt.Println(user2)

	dbUser, err := testUserDao.GetByLogin("test_login")
	test.TestForError(t, err)

	fmt.Println(dbUser)
}
